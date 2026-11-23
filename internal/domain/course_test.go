package domain_test

import (
	"errors"
	"math"
	"testing"

	"example.com/zoo/internal/domain"
)

func TestAddStudent(t *testing.T) {
	t.Parallel()
	// Arrange: допустимые и недопустимые данные.
	tests := []struct {
		name, id, student string
		want              error
	}{
		{"корректный", "s1", " Анна Иванова ", nil},
		{"другой ID с тем же именем", "s2", "Анна", nil},
		{"пустой ID", "", "Анна", domain.ErrInvalidID},
		{"пробел в ID", "s 1", "Анна", domain.ErrInvalidID},
		{"кириллица в ID", "студент", "Анна", domain.ErrInvalidID},
		{"пустое имя", "s1", "  ", domain.ErrEmptyName},
		{"дубликат", "existing", "Другое имя", domain.ErrDuplicate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			course := &domain.Course{}
			must(t, course.AddStudent("existing", "Анна"))
			// Act: добавляем студента.
			err := course.AddStudent(tt.id, tt.student)
			// Assert: ошибки не меняют список.
			if !errors.Is(err, tt.want) {
				t.Fatalf("ошибка = %v, ожидается %v", err, tt.want)
			}
			wantCount := 1
			if tt.want == nil {
				wantCount++
			}
			if len(course.Students()) != wantCount {
				t.Fatal("неверное число студентов")
			}
		})
	}
}

func TestAddCriterion(t *testing.T) {
	t.Parallel()
	// Arrange: границы и специальные числовые значения.
	tests := []struct {
		name, id, title string
		maximum, weight float64
		want            error
	}{
		{"корректный", "exam", "Экзамен", 10, 3, nil},
		{"пустой ID", "", "Экзамен", 10, 3, domain.ErrInvalidID},
		{"пустое название", "exam", " ", 10, 3, domain.ErrEmptyName},
		{"нулевой максимум", "exam", "Экзамен", 0, 3, domain.ErrInvalidMaximum},
		{"отрицательный максимум", "exam", "Экзамен", -1, 3, domain.ErrInvalidMaximum},
		{"NaN максимум", "exam", "Экзамен", math.NaN(), 3, domain.ErrInvalidMaximum},
		{"Inf максимум", "exam", "Экзамен", math.Inf(1), 3, domain.ErrInvalidMaximum},
		{"нулевой вес", "exam", "Экзамен", 10, 0, domain.ErrInvalidWeight},
		{"отрицательный вес", "exam", "Экзамен", 10, -1, domain.ErrInvalidWeight},
		{"NaN вес", "exam", "Экзамен", 10, math.NaN(), domain.ErrInvalidWeight},
		{"Inf вес", "exam", "Экзамен", 10, math.Inf(-1), domain.ErrInvalidWeight},
		{"дубликат", "existing", "Экзамен", 10, 3, domain.ErrDuplicate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			course := &domain.Course{}
			must(t, course.AddCriterion("existing", "Практика", 20, 1))
			// Act: добавляем критерий.
			err := course.AddCriterion(tt.id, tt.title, tt.maximum, tt.weight)
			// Assert: неверный критерий не сохраняется.
			if !errors.Is(err, tt.want) {
				t.Fatalf("ошибка = %v, ожидается %v", err, tt.want)
			}
			wantCount := 1
			if tt.want == nil {
				wantCount++
			}
			if len(course.Criteria()) != wantCount {
				t.Fatal("неверное число критериев")
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	t.Parallel()
	// Arrange: исходная оценка 5/10 и попытки её заменить.
	tests := []struct {
		name, student, criterion string
		score                    float64
		want                     error
	}{
		{"ноль", "s1", "exam", 0, nil},
		{"максимум", "s1", "exam", 10, nil},
		{"дробь", "s1", "exam", 8.5, nil},
		{"выше максимума", "s1", "exam", 12, domain.ErrInvalidScore},
		{"ниже нуля", "s1", "exam", -1, domain.ErrInvalidScore},
		{"NaN", "s1", "exam", math.NaN(), domain.ErrInvalidScore},
		{"Inf", "s1", "exam", math.Inf(1), domain.ErrInvalidScore},
		{"нет студента", "unknown", "exam", 8, domain.ErrStudentNotFound},
		{"нет критерия", "s1", "unknown", 8, domain.ErrCriterionNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			course := newCourse(t)
			must(t, course.Evaluate("s1", "exam", 5))
			// Act: фабрика проверяет балл и обе связи.
			err := course.Evaluate(tt.student, tt.criterion, tt.score)
			// Assert: неудачная попытка сохраняет прежнюю оценку.
			if !errors.Is(err, tt.want) {
				t.Fatalf("ошибка = %v, ожидается %v", err, tt.want)
			}
			wantScore := 5.0
			if tt.want == nil {
				wantScore = tt.score
			}
			got, exists := course.Evaluation("s1", "exam")
			if !exists || got.Score() != wantScore || got.StudentID() != "s1" || got.CriterionID() != "exam" {
				t.Fatalf("неверная оценка: %+v, наличие %v", got, exists)
			}
		})
	}
}

func TestCancelEvaluation(t *testing.T) {
	t.Parallel()
	// Arrange: отмена существующей и отсутствующих оценок.
	for _, tt := range []struct {
		name, student, criterion string
		graded                   bool
		want                     error
	}{
		{"отмена", "s1", "exam", true, nil},
		{"нет оценки", "s1", "exam", false, domain.ErrEvaluationNotFound},
		{"нет студента", "unknown", "exam", true, domain.ErrStudentNotFound},
		{"нет критерия", "s1", "unknown", true, domain.ErrCriterionNotFound},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			course := newCourse(t)
			if tt.graded {
				must(t, course.Evaluate("s1", "exam", 0))
			}
			// Act: отменяем оценку.
			err := course.CancelEvaluation(tt.student, tt.criterion)
			// Assert: отменённый ноль становится пропуском.
			if !errors.Is(err, tt.want) {
				t.Fatalf("ошибка = %v, ожидается %v", err, tt.want)
			}
			_, exists := course.Evaluation("s1", "exam")
			if exists != (tt.graded && tt.want != nil) {
				t.Fatal("неверное состояние оценки")
			}
		})
	}
}

func newCourse(t *testing.T) *domain.Course {
	t.Helper()
	course := &domain.Course{}
	must(t, course.AddStudent("s1", "Анна Иванова"))
	must(t, course.AddCriterion("exam", "Экзамен", 10, 3))
	must(t, course.AddCriterion("lab", "Практика", 20, 1))
	return course
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
