// Package domain содержит модель оценивания без ввода-вывода.
package domain

import "errors"

var (
	// ErrInvalidID означает недопустимый идентификатор.
	ErrInvalidID = errors.New("ID должен содержать только латинские буквы, цифры, _ и -")
	// ErrEmptyName означает пустое имя или название.
	ErrEmptyName = errors.New("имя или название не может быть пустым")
	// ErrDuplicate означает повторный идентификатор.
	ErrDuplicate = errors.New("ID уже существует")
	// ErrStudentNotFound означает отсутствие студента на курсе.
	ErrStudentNotFound = errors.New("студент не найден")
	// ErrCriterionNotFound означает отсутствие критерия в рубрикаторе.
	ErrCriterionNotFound = errors.New("критерий не найден")
	// ErrInvalidMaximum означает недопустимый максимум баллов.
	ErrInvalidMaximum = errors.New("максимум должен быть конечным положительным числом")
	// ErrInvalidWeight означает недопустимый вес критерия.
	ErrInvalidWeight = errors.New("вес должен быть конечным положительным числом")
	// ErrInvalidScore означает выход балла за допустимый диапазон.
	ErrInvalidScore = errors.New("балл должен быть конечным числом от 0 до максимума критерия")
	// ErrEvaluationNotFound означает отсутствие оценки.
	ErrEvaluationNotFound = errors.New("оценка не выставлена")
	// ErrEmptyRubric означает отсутствие критериев для расчёта.
	ErrEmptyRubric = errors.New("рубрикатор пуст")
	// ErrIncomplete означает неполный набор оценок для окончательного итога.
	ErrIncomplete = errors.New("выставлены не все оценки")
)
