package zoo

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
)

// Zoo — зоопарк: набор вольеров, адресуемых по названию.
type Zoo struct {
	name       string
	enclosures map[string]*Enclosure
}

// New создаёт пустой зоопарк.
func New(name string) *Zoo {
	return &Zoo{name: name, enclosures: make(map[string]*Enclosure)}
}

// Name возвращает название зоопарка.
func (z *Zoo) Name() string {
	return z.name
}

// AddEnclosure добавляет вольер в зоопарк.
func (z *Zoo) AddEnclosure(e *Enclosure) error {
	if _, ok := z.enclosures[e.Name()]; ok {
		return fmt.Errorf("%w: %s", ErrDuplicateEnclosure, e.Name())
	}
	z.enclosures[e.Name()] = e
	return nil
}

// Enclosure возвращает вольер по названию.
func (z *Zoo) Enclosure(name string) (*Enclosure, error) {
	e, ok := z.enclosures[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnclosureNotFound, name)
	}
	return e, nil
}

// Enclosures возвращает все вольеры, отсортированные по названию.
//
// maps.Values + slices.SortedFunc — современная (Go 1.23+) альтернатива
// ручному циклу сбора значений и сортировки.
func (z *Zoo) Enclosures() []*Enclosure {
	return slices.SortedFunc(maps.Values(z.enclosures), func(a, b *Enclosure) int {
		return cmp.Compare(a.Name(), b.Name())
	})
}
