package domain

import (
	"math"
	"strings"
)

// Student — студент курса с устойчивым ID.
type Student struct {
	id, name string
}

// ID возвращает идентификатор студента.
func (s Student) ID() string { return s.id }

// Name возвращает имя студента.
func (s Student) Name() string { return s.name }

// Criterion — критерий рубрикатора с устойчивым ID.
type Criterion struct {
	id, name        string
	maximum, weight float64
}

// ID возвращает идентификатор критерия.
func (c Criterion) ID() string { return c.id }

// Name возвращает название критерия.
func (c Criterion) Name() string { return c.name }

// Maximum возвращает максимально допустимый балл.
func (c Criterion) Maximum() float64 { return c.maximum }

// Weight возвращает относительный вес критерия.
func (c Criterion) Weight() float64 { return c.weight }

// Evaluation — неизменяемое значение оценки студента по критерию.
type Evaluation struct {
	studentID, criterionID string
	score                  float64
}

// StudentID возвращает идентификатор студента.
func (e Evaluation) StudentID() string { return e.studentID }

// CriterionID возвращает идентификатор критерия.
func (e Evaluation) CriterionID() string { return e.criterionID }

// Score возвращает выставленный балл.
func (e Evaluation) Score() float64 { return e.score }

func validateIdentity(id, name string) error {
	if id == "" || strings.ContainsFunc(id, func(r rune) bool {
		letter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		digit := r >= '0' && r <= '9'
		return !letter && !digit && r != '_' && r != '-'
	}) {
		return ErrInvalidID
	}
	if strings.TrimSpace(name) == "" {
		return ErrEmptyName
	}
	return nil
}

func finite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
