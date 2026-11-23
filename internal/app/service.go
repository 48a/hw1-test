// Package app содержит сценарии ПОС и порты внешних систем.
package app

import (
	"fmt"

	"example.com/zoo/internal/domain"
)

// Repository — порт хранения курса; чтение и запись передают независимые копии.
type Repository interface {
	// Load возвращает независимый снимок курса.
	Load() (*domain.Course, error)
	// Save сохраняет независимую копию курса.
	Save(course *domain.Course) error
}

// Exporter — порт выгрузки вычисленных итогов.
type Exporter interface {
	// Export выгружает готовые итоги в указанное место.
	Export(path string, results []Result) error
}

// Result — строка отчёта со студентом и его текущим итогом.
type Result struct {
	// Student — студент курса.
	Student domain.Student
	// Total — рассчитанный итог.
	Total domain.Total
}

// Summary — текущие итоги и статистика по всем студентам.
type Summary struct {
	// Results — строки, упорядоченные по ID студента.
	Results []Result
	// Average — средний процент по группе.
	Average float64
	// Minimum — наименьший процент в группе.
	Minimum float64
	// Maximum — наибольший процент в группе.
	Maximum float64
}

// Service — фасад сценариев для консоли.
type Service struct {
	repository Repository
	exporter   Exporter
	current    domain.Calculator
	final      domain.Calculator
}

// New создаёт фасад с зависимостями, переданными DI-контейнером.
func New(repository Repository, exporter Exporter, current, final domain.Calculator) *Service {
	return &Service{repository: repository, exporter: exporter, current: current, final: final}
}

// Snapshot возвращает независимый снимок курса для просмотра.
func (s *Service) Snapshot() (*domain.Course, error) {
	course, err := s.repository.Load()
	if err != nil {
		return nil, fmt.Errorf("загрузка курса: %w", err)
	}
	return course, nil
}

// AddStudent добавляет студента и сохраняет курс.
func (s *Service) AddStudent(id, name string) error {
	return s.change(func(course *domain.Course) error {
		return course.AddStudent(id, name)
	})
}

// AddCriterion добавляет критерий и сохраняет курс.
func (s *Service) AddCriterion(id, name string, maximum, weight float64) error {
	return s.change(func(course *domain.Course) error {
		return course.AddCriterion(id, name, maximum, weight)
	})
}

// Grade выставляет или исправляет оценку.
func (s *Service) Grade(studentID, criterionID string, score float64) error {
	return s.change(func(course *domain.Course) error {
		return course.Evaluate(studentID, criterionID, score)
	})
}

// Ungrade отменяет оценку.
func (s *Service) Ungrade(studentID, criterionID string) error {
	return s.change(func(course *domain.Course) error {
		return course.CancelEvaluation(studentID, criterionID)
	})
}

// Total вычисляет текущий или окончательный итог студента.
func (s *Service) Total(studentID string, final bool) (domain.Total, error) {
	course, err := s.Snapshot()
	if err != nil {
		return domain.Total{}, err
	}
	calculator := s.current
	if final {
		calculator = s.final
	}
	return calculator.Calculate(course, studentID)
}

// Summary вычисляет текущие итоги и статистику из одного снимка курса.
func (s *Service) Summary() (Summary, error) {
	course, err := s.Snapshot()
	if err != nil {
		return Summary{}, err
	}
	var summary Summary
	for _, student := range course.Students() {
		total, err := s.current.Calculate(course, student.ID())
		if err != nil {
			return Summary{}, fmt.Errorf("студент %s: %w", student.ID(), err)
		}
		summary.Results = append(summary.Results, Result{Student: student, Total: total})
		summary.Average += total.Percent
		if len(summary.Results) == 1 {
			summary.Minimum = total.Percent
		}
		summary.Minimum = min(summary.Minimum, total.Percent)
		summary.Maximum = max(summary.Maximum, total.Percent)
	}
	if len(summary.Results) > 0 {
		summary.Average /= float64(len(summary.Results))
	}
	return summary, nil
}

// Export выгружает текущие итоги через выходной порт.
func (s *Service) Export(path string) error {
	summary, err := s.Summary()
	if err != nil {
		return err
	}
	if err := s.exporter.Export(path, summary.Results); err != nil {
		return fmt.Errorf("экспорт: %w", err)
	}
	return nil
}

func (s *Service) change(change func(*domain.Course) error) error {
	course, err := s.Snapshot()
	if err != nil {
		return err
	}
	if err := change(course); err != nil {
		return err
	}
	if err := s.repository.Save(course); err != nil {
		return fmt.Errorf("сохранение курса: %w", err)
	}
	return nil
}
