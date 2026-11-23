// Package memory хранит курс в памяти одного процесса.
package memory

import "example.com/zoo/internal/domain"

// Repository — хранилище курса для однопоточного REPL.
// Нулевое значение содержит пустой курс.
type Repository struct {
	course domain.Course
}

// Load возвращает независимую копию курса.
func (r *Repository) Load() (*domain.Course, error) {
	return r.course.Clone(), nil
}

// Save сохраняет независимую копию курса.
func (r *Repository) Save(course *domain.Course) error {
	r.course = *course.Clone()
	return nil
}
