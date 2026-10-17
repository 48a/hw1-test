package zoo

import (
	"fmt"
	"slices"
)

// Enclosure — вольер с ограниченной вместимостью.
//
// Поля не экспортированы: состояние меняется только через методы,
// так инварианты (вместимость, уникальность) проще удержать.
type Enclosure struct {
	name     string
	capacity int
	animals  []Animal
}

// NewEnclosure создаёт вольер заданной вместимости.
func NewEnclosure(name string, capacity int) (*Enclosure, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: пустое название", ErrInvalidEnclosure)
	}
	if capacity <= 0 {
		return nil, fmt.Errorf("%w: вместимость должна быть положительной", ErrInvalidEnclosure)
	}
	return &Enclosure{name: name, capacity: capacity}, nil
}

// Name возвращает название вольера.
func (e *Enclosure) Name() string {
	return e.name
}

// Capacity возвращает вместимость вольера.
func (e *Enclosure) Capacity() int {
	return e.capacity
}

// HasSpace сообщает, есть ли в вольере свободные места.
func (e *Enclosure) HasSpace() bool {
	return len(e.animals) < e.capacity
}

// Add помещает животное в вольер.
func (e *Enclosure) Add(a Animal) error {
	if !e.HasSpace() {
		return fmt.Errorf("%w: %s", ErrEnclosureFull, e.name)
	}
	e.animals = append(e.animals, a)
	return nil
}

// Remove забирает животное из вольера по имени.
func (e *Enclosure) Remove(name string) (Animal, error) {
	i := slices.IndexFunc(e.animals, func(a Animal) bool { return a.Name == name })
	if i < 0 {
		return Animal{}, fmt.Errorf("%w: %s", ErrAnimalNotFound, name)
	}
	a := e.animals[i]
	e.animals = slices.Delete(e.animals, i, i+1)
	return a, nil
}

// Animals возвращает копию списка животных вольера.
// Копия — чтобы вызывающий код не мог изменить состояние вольера в обход методов.
func (e *Enclosure) Animals() []Animal {
	return slices.Clone(e.animals)
}
