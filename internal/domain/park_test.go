package domain_test

import (
	"errors"
	"math"
	"slices"
	"strconv"
	"testing"

	"example.com/zoo/internal/domain"
)

func newPark(t *testing.T, service domain.ServiceCenter) *domain.Park {
	t.Helper()
	park, err := domain.NewPark(service)
	if err != nil {
		t.Fatal(err)
	}
	return park
}

func bicycle(t *testing.T, number string, ease int, operational bool) domain.Bicycle {
	t.Helper()
	transport, err := domain.NewBicycle("Велосипед", number, ease, operational)
	if err != nil {
		t.Fatal(err)
	}
	return transport
}

func TestNewParkRequiresService(t *testing.T) {
	t.Parallel()
	// Arrange: отсутствующая зависимость техосмотра.
	var service domain.ServiceCenter
	// Act: создаём парк без сервиса.
	park, err := domain.NewPark(service)
	// Assert: конструктор отказывает, не создавая неработоспособный парк.
	if park != nil || !errors.Is(err, domain.ErrServiceRequired) {
		t.Fatalf("NewPark(nil) = %v, %v", park, err)
	}
}

func TestAddTransportInspection(t *testing.T) {
	t.Parallel()
	// Arrange: решение мока намеренно не совпадает с исправностью.
	tests := []struct {
		name        string
		operational bool
		accepted    bool
		wantErr     error
		wantTotal   int
	}{
		{"сервис принимает неисправный", false, true, nil, 1},
		{"сервис отклоняет исправный", true, false, domain.ErrInspectionRejected, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			transport := bicycle(t, "1", 6, tt.operational)
			service := NewMockServiceCenter(t)
			service.EXPECT().Inspect(transport).Return(tt.accepted).Once()
			park := newPark(t, service)
			// Act: пытаемся принять транспорт.
			err := park.AddTransport(transport)
			// Assert: баланс следует решению сервиса; мок проверяет ровно один вызов.
			if !errors.Is(err, tt.wantErr) || park.Report().Total != tt.wantTotal {
				t.Fatalf("ошибка = %v, отчёт = %+v", err, park.Report())
			}
		})
	}
}

func TestBasicServiceCenter(t *testing.T) {
	t.Parallel()
	// Arrange: два технических состояния.
	for _, operational := range []bool{false, true} {
		transport := bicycle(t, "1", 6, operational)
		// Act: запускаем настоящую проверку.
		accepted := (domain.BasicServiceCenter{}).Inspect(transport)
		// Assert: принимается только исправный транспорт.
		if accepted != operational {
			t.Errorf("исправность = %v, решение = %v", operational, accepted)
		}
	}
}

func TestRejectionDoesNotReserveNumber(t *testing.T) {
	t.Parallel()
	// Arrange: первый техосмотр отклоняет транспорт, повторный принимает.
	transport := bicycle(t, "1", 6, true)
	service := NewMockServiceCenter(t)
	service.EXPECT().Inspect(transport).Return(false).Once()
	service.EXPECT().Inspect(transport).Return(true).Once()
	park := newPark(t, service)
	// Act: повторяем приёмку с тем же номером.
	first := park.AddTransport(transport)
	second := park.AddTransport(transport)
	// Assert: отказ не занимает номер и не добавляет объект в отчёты.
	if !errors.Is(first, domain.ErrInspectionRejected) || second != nil || park.Report().Total != 1 {
		t.Fatalf("первый результат = %v, второй = %v, отчёт = %+v", first, second, park.Report())
	}
}

func TestDuplicateNumbers(t *testing.T) {
	t.Parallel()
	// Arrange: одинаковые номера у разных групп объектов, включая ведущие нули.
	for _, firstTransport := range []bool{false, true} {
		t.Run(map[bool]string{false: "сначала вещь", true: "сначала транспорт"}[firstTransport], func(t *testing.T) {
			t.Parallel()
			transport := bicycle(t, "001", 6, true)
			helmet, err := domain.NewHelmet("Шлем", "1")
			if err != nil {
				t.Fatal(err)
			}
			service := NewMockServiceCenter(t)
			park := newPark(t, service)
			if firstTransport {
				service.EXPECT().Inspect(transport).Return(true).Once()
				if err := park.AddTransport(transport); err != nil {
					t.Fatal(err)
				}
			} else if err := park.AddItem(helmet); err != nil {
				t.Fatal(err)
			}
			// Act: повторяем приёмку в обеих группах.
			errorsGot := []error{park.AddTransport(transport), park.AddItem(helmet)}
			// Assert: дубли отклонены до техосмотра; баланс не изменился.
			for _, err := range errorsGot {
				if !errors.Is(err, domain.ErrDuplicateNumber) {
					t.Errorf("ошибка = %v", err)
				}
			}
			if park.Report().Total != 1 {
				t.Errorf("баланс изменился: %+v", park.Report())
			}
		})
	}
}

func TestInvalidAdmission(t *testing.T) {
	t.Parallel()
	// Arrange: нулевые значения не должны попасть на баланс или в техосмотр.
	tests := []struct {
		name  string
		add   func(*domain.Park) error
		error error
	}{
		{"nil транспорт", func(p *domain.Park) error { return p.AddTransport(nil) }, domain.ErrInvalidAsset},
		{"нулевой самокат", func(p *domain.Park) error { return p.AddTransport(domain.Scooter{}) }, domain.ErrInvalidAsset},
		{"nil вещь", func(p *domain.Park) error { return p.AddItem(nil) }, domain.ErrInvalidAsset},
		{"нулевой шлем", func(p *domain.Park) error { return p.AddItem(domain.Helmet{}) }, domain.ErrInvalidAsset},
		{"обход техосмотра", func(p *domain.Park) error { return p.AddItem(bicycle(t, "1", 6, true)) }, domain.ErrTransportAsItem},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			park := newPark(t, NewMockServiceCenter(t))
			// Act: пытаемся принять объект.
			err := tt.add(park)
			// Assert: баланс пуст, мок не получал вызовов.
			if !errors.Is(err, tt.error) || park.Report().Total != 0 {
				t.Fatalf("ошибка = %v, отчёт = %+v", err, park.Report())
			}
		})
	}
}

func TestBeginners(t *testing.T) {
	t.Parallel()
	// Arrange: принятый транспорт на границах правила и отклонённый простой.
	park := newPark(t, domain.BasicServiceCenter{})
	for i, ease := range []int{1, 5, 6, 10} {
		transport := bicycle(t, strconv.Itoa(i+1), ease, true)
		if err := park.AddTransport(transport); err != nil {
			t.Fatal(err)
		}
	}
	if err := park.AddTransport(bicycle(t, "5", 10, false)); !errors.Is(err, domain.ErrInspectionRejected) {
		t.Fatal(err)
	}
	helmet, err := domain.NewHelmet("Шлем", "6")
	if err != nil {
		t.Fatal(err)
	}
	if err := park.AddItem(helmet); err != nil {
		t.Fatal(err)
	}
	// Act: запрашиваем подборку новичкам.
	got := park.Beginners()
	// Assert: включены только принятые 6 и 10, с сохранением порядка приёмки.
	if len(got) != 2 || got[0].BeginnerEase() != 6 || got[1].BeginnerEase() != 10 {
		t.Fatalf("подборка = %v", got)
	}
	got[0] = nil
	if park.Beginners()[0] == nil {
		t.Error("изменение результата повредило парк")
	}
}

func TestParkReport(t *testing.T) {
	t.Parallel()
	// Arrange: сценарии пустого и заполненного баланса.
	for _, populated := range []bool{false, true} {
		t.Run(map[bool]string{false: "пустой", true: "полный"}[populated], func(t *testing.T) {
			t.Parallel()
			park := newPark(t, domain.BasicServiceCenter{})
			var want domain.Report
			if populated {
				populatePark(t, park)
				want = domain.Report{Transport: 3, Items: 3, Total: 6, Energy: 6.75}
			}
			// Act: формируем отчёт и получаем весь инвентарь.
			got, inventory := park.Report(), park.Inventory()
			// Assert: учтены все виды, в том числе энергия зарядного шкафа.
			if got != want || len(inventory) != want.Total {
				t.Fatalf("отчёт = %+v, ожидался %+v; инвентарь = %d", got, want, len(inventory))
			}
			if populated {
				inventory[0] = nil
				if park.Inventory()[0] == nil {
					t.Error("изменение списка повредило парк")
				}
			}
		})
	}
}

type electricTricycle struct {
	domain.Bicycle
	ease  int
	daily float64
}

// BeginnerEase возвращает простоту нового вида транспорта.
func (v electricTricycle) BeginnerEase() int { return v.ease }

// DailyEnergy возвращает расход нового вида транспорта.
func (v electricTricycle) DailyEnergy() float64 { return v.daily }

func TestNewTransportType(t *testing.T) {
	t.Parallel()
	// Arrange: тип из другого пакета, о котором парк ничего не знает.
	tests := []struct {
		name    string
		ease    int
		daily   float64
		wantErr error
	}{
		{"корректный трицикл", 7, 4.5, nil},
		{"слишком малая простота", 0, 4.5, domain.ErrInvalidEase},
		{"слишком большая простота", 11, 4.5, domain.ErrInvalidEase},
		{"отрицательный расход", 7, -1, domain.ErrInvalidEnergy},
		{"расход NaN", 7, math.NaN(), domain.ErrInvalidEnergy},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			transport := electricTricycle{Bicycle: bicycle(t, "1", 6, true), ease: tt.ease, daily: tt.daily}
			service := NewMockServiceCenter(t)
			park := newPark(t, service)
			if tt.wantErr == nil {
				service.EXPECT().Inspect(transport).Return(true).Once()
			}
			// Act: принимаем новый тип через общий интерфейс.
			err := park.AddTransport(transport)
			// Assert: контракты проверены до техосмотра, отчёты учитывают новый вид.
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ошибка = %v, ожидалась %v", err, tt.wantErr)
			}
			var want domain.Report
			if tt.wantErr == nil {
				want = domain.Report{Transport: 1, Total: 1, Energy: tt.daily}
			}
			if park.Report() != want || len(park.Beginners()) != want.Transport {
				t.Fatalf("отчёт = %+v, подборка = %v", park.Report(), park.Beginners())
			}
		})
	}
}

func populatePark(t *testing.T, park *domain.Park) {
	t.Helper()
	scooter, err := domain.NewScooter("Самокат", "1", 6, 1.5, true)
	if err != nil {
		t.Fatal(err)
	}
	electric, err := domain.NewElectricBicycle("Электровелосипед", "2", 5, 2.25, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, transport := range []domain.Transport{scooter, electric, bicycle(t, "3", 8, true)} {
		if err := park.AddTransport(transport); err != nil {
			t.Fatal(err)
		}
	}
	helmet, err := domain.NewHelmet("Шлем", "4")
	if err != nil {
		t.Fatal(err)
	}
	dock, err := domain.NewDockStation("Док", "5")
	if err != nil {
		t.Fatal(err)
	}
	cabinet, err := domain.NewChargingCabinet("Шкаф", "6", 3)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range []domain.Inventory{helmet, dock, cabinet} {
		if err := park.AddItem(item); err != nil {
			t.Fatal(err)
		}
	}
	wantNumbers := []string{"1", "2", "3", "4", "5", "6"}
	var numbers []string
	for _, item := range park.Inventory() {
		numbers = append(numbers, item.InventoryNumber())
	}
	if !slices.Equal(numbers, wantNumbers) {
		t.Fatalf("номера = %v", numbers)
	}
}
