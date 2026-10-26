package domain

import "math"

// EnergyConsumer описывает потребителя электроэнергии.
type EnergyConsumer interface {
	// DailyEnergy возвращает конечный неотрицательный расход в кВт·ч.
	DailyEnergy() float64
}

type energy struct{ daily float64 }

// DailyEnergy возвращает суточный расход в кВт·ч.
func (e energy) DailyEnergy() float64 { return e.daily }

func newEnergy(daily float64) (energy, error) {
	if daily < 0 || math.IsNaN(daily) || math.IsInf(daily, 0) {
		return energy{}, ErrInvalidEnergy
	}
	return energy{daily: daily}, nil
}
