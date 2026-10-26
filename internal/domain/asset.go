package domain

import "strings"

// Inventory описывает объект на балансе компании.
type Inventory interface {
	// Name возвращает название объекта.
	Name() string
	// InventoryNumber возвращает положительный десятичный номер.
	InventoryNumber() string
}

type asset struct {
	name   string
	number string
}

// Name возвращает название объекта.
func (a asset) Name() string { return a.name }

// InventoryNumber возвращает инвентарный номер без ведущих нулей.
func (a asset) InventoryNumber() string { return a.number }

func newAsset(name, number string) (asset, error) {
	name, number = strings.TrimSpace(name), strings.TrimSpace(number)
	if name == "" || number == "" {
		return asset{}, ErrInvalidAsset
	}
	for _, digit := range number {
		if digit < '0' || digit > '9' {
			return asset{}, ErrInvalidAsset
		}
	}
	number = strings.TrimLeft(number, "0")
	if number == "" {
		return asset{}, ErrInvalidAsset
	}
	return asset{name: name, number: number}, nil
}
