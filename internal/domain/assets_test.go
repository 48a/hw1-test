package domain_test

import (
	"errors"
	"math"
	"testing"

	"example.com/zoo/internal/domain"
)

func TestAssetConstructors(t *testing.T) {
	t.Parallel()
	// Arrange: все шесть видов с общими полями и ожидаемыми свойствами.
	tests := []struct {
		name      string
		create    func(string, string) (domain.Inventory, error)
		transport bool
		energy    bool
	}{
		{"самокат", func(n, id string) (domain.Inventory, error) { return domain.NewScooter(n, id, 6, 1.5, true) }, true, true},
		{"велосипед", func(n, id string) (domain.Inventory, error) { return domain.NewBicycle(n, id, 6, true) }, true, false},
		{"электровелосипед", func(n, id string) (domain.Inventory, error) { return domain.NewElectricBicycle(n, id, 6, 1.5, true) }, true, true},
		{"шлем", func(n, id string) (domain.Inventory, error) { return domain.NewHelmet(n, id) }, false, false},
		{"док-станция", func(n, id string) (domain.Inventory, error) { return domain.NewDockStation(n, id) }, false, false},
		{"зарядный шкаф", func(n, id string) (domain.Inventory, error) { return domain.NewChargingCabinet(n, id, 1.5) }, false, true},
	}
	inputs := []struct {
		name, number string
		wantErr      error
	}{
		{"  Объект с пробелами  ", "  00123  ", nil},
		{"", "123", domain.ErrInvalidAsset},
		{"   ", "123", domain.ErrInvalidAsset},
		{"Название", "", domain.ErrInvalidAsset},
		{"Название", "000", domain.ErrInvalidAsset},
		{"Название", "-1", domain.ErrInvalidAsset},
		{"Название", "+1", domain.ErrInvalidAsset},
		{"Название", "1.5", domain.ErrInvalidAsset},
		{"Название", "abc", domain.ErrInvalidAsset},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for _, input := range inputs {
				t.Run(input.name+"/"+input.number, func(t *testing.T) {
					t.Parallel()
					// Act: создаём объект из пользовательских данных.
					item, err := tt.create(input.name, input.number)
					// Assert: ошибка соответствует причине, успешный объект нормализован.
					if !errors.Is(err, input.wantErr) {
						t.Fatalf("ошибка = %v, ожидалась %v", err, input.wantErr)
					}
					if input.wantErr != nil {
						return
					}
					assertProperties(t, item, tt.transport, tt.energy)
				})
			}
		})
	}
}

func assertProperties(t *testing.T, item domain.Inventory, wantTransport, wantEnergy bool) {
	t.Helper()
	if item.Name() != "Объект с пробелами" || item.InventoryNumber() != "123" {
		t.Fatalf("поля объекта: %q [%s]", item.Name(), item.InventoryNumber())
	}
	transport, isTransport := item.(domain.Transport)
	if isTransport != wantTransport {
		t.Fatalf("Transport = %v, ожидалось %v", isTransport, wantTransport)
	}
	if isTransport && (transport.BeginnerEase() != 6 || !transport.IsOperational()) {
		t.Error("не сохранены свойства транспорта")
	}
	consumer, consumes := item.(domain.EnergyConsumer)
	if consumes != wantEnergy {
		t.Fatalf("EnergyConsumer = %v, ожидалось %v", consumes, wantEnergy)
	}
	if consumes && consumer.DailyEnergy() != 1.5 {
		t.Errorf("энергия = %v, ожидалось 1.5", consumer.DailyEnergy())
	}
}

func TestTransportEase(t *testing.T) {
	t.Parallel()
	// Arrange: границы допустимой простоты для каждого вида транспорта.
	constructors := []struct {
		name   string
		create func(int) (domain.Transport, error)
	}{
		{"самокат", func(ease int) (domain.Transport, error) {
			return domain.NewScooter("Самокат", "1", ease, 0, false)
		}},
		{"велосипед", func(ease int) (domain.Transport, error) {
			return domain.NewBicycle("Велосипед", "1", ease, false)
		}},
		{"электровелосипед", func(ease int) (domain.Transport, error) {
			return domain.NewElectricBicycle("Электровелосипед", "1", ease, 0, false)
		}},
	}
	for _, constructor := range constructors {
		t.Run(constructor.name, func(t *testing.T) {
			t.Parallel()
			for _, ease := range []int{-1, 0, 1, 10, 11} {
				// Act: задаём простоту, в том числе значения за границами.
				transport, err := constructor.create(ease)
				// Assert: 1 и 10 допустимы, остальные значения отклонены.
				if ease < 1 || ease > 10 {
					if !errors.Is(err, domain.ErrInvalidEase) {
						t.Errorf("простота %d: ошибка = %v", ease, err)
					}
					continue
				}
				if err != nil {
					t.Fatal(err)
				}
				if transport.BeginnerEase() != ease || transport.IsOperational() {
					t.Error("не сохранены простота или неисправное состояние")
				}
			}
		})
	}
}

func TestEnergyValidation(t *testing.T) {
	t.Parallel()
	// Arrange: потребители энергии и допустимые/недопустимые расходы.
	constructors := []struct {
		name   string
		create func(float64) (domain.EnergyConsumer, error)
	}{
		{"самокат", func(daily float64) (domain.EnergyConsumer, error) {
			return domain.NewScooter("Самокат", "1", 6, daily, true)
		}},
		{"электровелосипед", func(daily float64) (domain.EnergyConsumer, error) {
			return domain.NewElectricBicycle("Велосипед", "1", 6, daily, true)
		}},
		{"шкаф", func(daily float64) (domain.EnergyConsumer, error) {
			return domain.NewChargingCabinet("Шкаф", "1", daily)
		}},
	}
	for _, constructor := range constructors {
		t.Run(constructor.name, func(t *testing.T) {
			t.Parallel()
			for _, daily := range []float64{0, 2.5, -1, math.NaN(), math.Inf(1), math.Inf(-1)} {
				// Act: создаём потребителя с заданным расходом.
				consumer, err := constructor.create(daily)
				// Assert: нечисловые, бесконечные и отрицательные расходы отклонены.
				if daily < 0 || math.IsNaN(daily) || math.IsInf(daily, 0) {
					if !errors.Is(err, domain.ErrInvalidEnergy) {
						t.Errorf("расход %v: ошибка = %v", daily, err)
					}
					continue
				}
				if err != nil {
					t.Fatal(err)
				}
				if consumer.DailyEnergy() != daily {
					t.Errorf("расход = %v, ожидалось %v", consumer.DailyEnergy(), daily)
				}
			}
		})
	}
}
