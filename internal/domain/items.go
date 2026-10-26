package domain

// Helmet представляет шлем.
type Helmet struct{ asset }

// NewHelmet создаёт шлем с проверкой данных.
func NewHelmet(name, number string) (Helmet, error) {
	a, err := newAsset(name, number)
	if err != nil {
		return Helmet{}, err
	}
	return Helmet{asset: a}, nil
}

// DockStation представляет пассивную док-станцию без зарядки.
type DockStation struct{ asset }

// NewDockStation создаёт док-станцию с проверкой данных.
func NewDockStation(name, number string) (DockStation, error) {
	a, err := newAsset(name, number)
	if err != nil {
		return DockStation{}, err
	}
	return DockStation{asset: a}, nil
}

// ChargingCabinet представляет зарядный шкаф.
type ChargingCabinet struct {
	asset
	energy
}

// NewChargingCabinet создаёт зарядный шкаф с проверкой данных.
func NewChargingCabinet(name, number string, daily float64) (ChargingCabinet, error) {
	a, err := newAsset(name, number)
	if err != nil {
		return ChargingCabinet{}, err
	}
	e, err := newEnergy(daily)
	if err != nil {
		return ChargingCabinet{}, err
	}
	return ChargingCabinet{asset: a, energy: e}, nil
}
