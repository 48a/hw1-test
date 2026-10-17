// Package zoo содержит предметную область зоопарка:
// животных, вольеры и сам зоопарк.
//
// Пакет ничего не знает о консоли: ввод-вывод — ответственность internal/cli.
package zoo

import "errors"

// Ошибки предметной области.
// Сравнивайте их через errors.Is, не по тексту сообщения.
var (
	ErrInvalidAnimal      = errors.New("некорректные данные животного")
	ErrInvalidEnclosure   = errors.New("некорректные данные вольера")
	ErrEnclosureNotFound  = errors.New("вольер не найден")
	ErrEnclosureFull      = errors.New("вольер переполнен")
	ErrDuplicateEnclosure = errors.New("вольер уже существует")
	ErrAnimalNotFound     = errors.New("животное не найдено")
)
