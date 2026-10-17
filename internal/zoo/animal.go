package zoo

import "fmt"

// Animal — животное в зоопарке.
type Animal struct {
	Name    string
	Species string
	Age     int
}

// NewAnimal создаёт животное, проверяя корректность полей.
func NewAnimal(name, species string, age int) (Animal, error) {
	if name == "" {
		return Animal{}, fmt.Errorf("%w: пустое имя", ErrInvalidAnimal)
	}
	if species == "" {
		return Animal{}, fmt.Errorf("%w: пустой вид", ErrInvalidAnimal)
	}
	if age < 0 {
		return Animal{}, fmt.Errorf("%w: возраст не может быть отрицательным", ErrInvalidAnimal)
	}
	return Animal{Name: name, Species: species, Age: age}, nil
}

// String возвращает краткое описание животного.
func (a Animal) String() string {
	return fmt.Sprintf("%s (%s, %d лет)", a.Name, a.Species, a.Age)
}
