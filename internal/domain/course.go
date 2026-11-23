package domain

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type evaluationKey struct {
	studentID, criterionID string
}

// Course — корень агрегата: студенты, рубрикатор и оценки одного курса.
// Нулевое значение готово к работе; все изменения проходят через методы.
type Course struct {
	students    map[string]Student
	criteria    map[string]Criterion
	evaluations map[evaluationKey]Evaluation
}

// AddStudent добавляет студента с уникальным ID.
func (c *Course) AddStudent(id, name string) error {
	if err := validateIdentity(id, name); err != nil {
		return err
	}
	if _, exists := c.students[id]; exists {
		return fmt.Errorf("%w: студент %s", ErrDuplicate, id)
	}
	if c.students == nil {
		c.students = make(map[string]Student)
	}
	c.students[id] = Student{id: id, name: strings.TrimSpace(name)}
	return nil
}

// AddCriterion добавляет критерий с положительными максимумом и весом.
func (c *Course) AddCriterion(id, name string, maximum, weight float64) error {
	if err := validateIdentity(id, name); err != nil {
		return err
	}
	if !finite(maximum) || maximum <= 0 {
		return ErrInvalidMaximum
	}
	if !finite(weight) || weight <= 0 {
		return ErrInvalidWeight
	}
	if _, exists := c.criteria[id]; exists {
		return fmt.Errorf("%w: критерий %s", ErrDuplicate, id)
	}
	if c.criteria == nil {
		c.criteria = make(map[string]Criterion)
	}
	c.criteria[id] = Criterion{id: id, name: strings.TrimSpace(name), maximum: maximum, weight: weight}
	return nil
}

// Evaluate создаёт оценку после проверки связей и балла, заменяя прежнее значение.
// Это единственная фабрика оценок внутри агрегата.
func (c *Course) Evaluate(studentID, criterionID string, score float64) error {
	criterion, err := c.evaluationCriterion(studentID, criterionID)
	if err != nil {
		return err
	}
	if !finite(score) || score < 0 || score > criterion.maximum {
		return fmt.Errorf("%w: максимум %g", ErrInvalidScore, criterion.maximum)
	}
	if c.evaluations == nil {
		c.evaluations = make(map[evaluationKey]Evaluation)
	}
	c.evaluations[evaluationKey{studentID, criterionID}] = Evaluation{studentID, criterionID, score}
	return nil
}

// CancelEvaluation отменяет существующую оценку.
func (c *Course) CancelEvaluation(studentID, criterionID string) error {
	if _, err := c.evaluationCriterion(studentID, criterionID); err != nil {
		return err
	}
	key := evaluationKey{studentID, criterionID}
	if _, exists := c.evaluations[key]; !exists {
		return ErrEvaluationNotFound
	}
	delete(c.evaluations, key)
	return nil
}

// Student находит студента по ID.
func (c *Course) Student(id string) (Student, error) {
	student, exists := c.students[id]
	if !exists {
		return Student{}, fmt.Errorf("%w: %s", ErrStudentNotFound, id)
	}
	return student, nil
}

// Students возвращает независимый список студентов по ID.
func (c *Course) Students() []Student {
	return slices.SortedFunc(maps.Values(c.students), func(a, b Student) int {
		return cmp.Compare(a.id, b.id)
	})
}

// Criteria возвращает независимый список критериев по ID.
func (c *Course) Criteria() []Criterion {
	return slices.SortedFunc(maps.Values(c.criteria), func(a, b Criterion) int {
		return cmp.Compare(a.id, b.id)
	})
}

// Evaluation возвращает оценку и признак её наличия, отличая пропуск от нуля.
func (c *Course) Evaluation(studentID, criterionID string) (Evaluation, bool) {
	evaluation, exists := c.evaluations[evaluationKey{studentID, criterionID}]
	return evaluation, exists
}

// Clone возвращает независимую копию агрегата для репозитория.
func (c *Course) Clone() *Course {
	return &Course{
		students: maps.Clone(c.students), criteria: maps.Clone(c.criteria),
		evaluations: maps.Clone(c.evaluations),
	}
}

func (c *Course) evaluationCriterion(studentID, criterionID string) (Criterion, error) {
	if _, err := c.Student(studentID); err != nil {
		return Criterion{}, err
	}
	criterion, exists := c.criteria[criterionID]
	if !exists {
		return Criterion{}, fmt.Errorf("%w: %s", ErrCriterionNotFound, criterionID)
	}
	return criterion, nil
}
