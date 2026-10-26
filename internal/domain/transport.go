package domain

const (
	minimumEase  = 1
	maximumEase  = 10
	beginnerEase = 6
)

// Transport описывает транспорт для техосмотра и подбора новичкам.
type Transport interface {
	Inventory
	// BeginnerEase возвращает простоту управления от 1 до 10.
	BeginnerEase() int
	// IsOperational сообщает об исправности транспорта.
	IsOperational() bool
}

type vehicle struct {
	asset
	ease        int
	operational bool
}

// BeginnerEase возвращает простоту управления от 1 до 10.
func (v vehicle) BeginnerEase() int { return v.ease }

// IsOperational сообщает об исправности транспорта.
func (v vehicle) IsOperational() bool { return v.operational }

func newVehicle(name, number string, ease int, operational bool) (vehicle, error) {
	a, err := newAsset(name, number)
	if err != nil {
		return vehicle{}, err
	}
	if ease < minimumEase || ease > maximumEase {
		return vehicle{}, ErrInvalidEase
	}
	return vehicle{asset: a, ease: ease, operational: operational}, nil
}

// Scooter представляет электросамокат.
type Scooter struct {
	vehicle
	energy
}

// NewScooter создаёт электросамокат с проверкой данных.
func NewScooter(name, number string, ease int, daily float64, operational bool) (Scooter, error) {
	v, err := newVehicle(name, number, ease, operational)
	if err != nil {
		return Scooter{}, err
	}
	e, err := newEnergy(daily)
	if err != nil {
		return Scooter{}, err
	}
	return Scooter{vehicle: v, energy: e}, nil
}

// Bicycle представляет велосипед без электропривода.
type Bicycle struct{ vehicle }

// NewBicycle создаёт велосипед с проверкой данных.
func NewBicycle(name, number string, ease int, operational bool) (Bicycle, error) {
	v, err := newVehicle(name, number, ease, operational)
	if err != nil {
		return Bicycle{}, err
	}
	return Bicycle{vehicle: v}, nil
}

// ElectricBicycle представляет электровелосипед.
type ElectricBicycle struct {
	vehicle
	energy
}

// NewElectricBicycle создаёт электровелосипед с проверкой данных.
func NewElectricBicycle(name, number string, ease int, daily float64, operational bool) (ElectricBicycle, error) {
	v, err := newVehicle(name, number, ease, operational)
	if err != nil {
		return ElectricBicycle{}, err
	}
	e, err := newEnergy(daily)
	if err != nil {
		return ElectricBicycle{}, err
	}
	return ElectricBicycle{vehicle: v, energy: e}, nil
}
