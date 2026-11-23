package domain_test

import (
	"errors"
	"math"
	"testing"

	"example.com/zoo/internal/domain"
)

func TestCalculate(t *testing.T) {
	t.Parallel()
	// Arrange: разные веса, пропуски и режимы расчёта.
	tests := []struct {
		name    string
		weights [2]float64
		scores  [2]*float64
		policy  domain.MissingPolicy
		want    float64
		graded  int
		wantErr error
	}{
		{"разные максимумы и веса", [2]float64{3, 1}, [2]*float64{new(8.0), new(20.0)}, domain.RequireAll{}, 85, 2, nil},
		{"дробные баллы", [2]float64{3, 1}, [2]*float64{new(7.5), new(5.0)}, nil, 62.5, 2, nil},
		{"дробные веса", [2]float64{0.75, 0.25}, [2]*float64{new(8.0), new(20.0)}, nil, 85, 2, nil},
		{"пропуск остаётся в знаменателе", [2]float64{3, 1}, [2]*float64{new(8.0), nil}, domain.MissingAsZero{}, 60, 1, nil},
		{"окончательный неполный", [2]float64{3, 1}, [2]*float64{new(8.0), nil}, domain.RequireAll{}, 0, 0, domain.ErrIncomplete},
		{"нет оценок", [2]float64{3, 1}, [2]*float64{}, nil, 0, 0, nil},
		{"ноль является оценкой", [2]float64{3, 1}, [2]*float64{new(0.0), new(0.0)}, domain.RequireAll{}, 0, 2, nil},
		{"большие веса", [2]float64{math.MaxFloat64, math.MaxFloat64}, [2]*float64{new(5.0), new(20.0)}, nil, 75, 2, nil},
		{"малые веса", [2]float64{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64}, [2]*float64{new(5.0), new(20.0)}, nil, 75, 2, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			course := &domain.Course{}
			must(t, course.AddStudent("s1", "Анна"))
			for i, id := range []string{"exam", "lab"} {
				must(t, course.AddCriterion(id, id, float64((i+1)*10), tt.weights[i]))
				if tt.scores[i] != nil {
					must(t, course.Evaluate("s1", id, *tt.scores[i]))
				}
			}
			// Act: доменный сервис получает текущий агрегат.
			got, err := domain.NewCalculator(tt.policy).Calculate(course, "s1")
			// Assert: нормировка не зависит от суммы весов и сохраняет пропуски.
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ошибка = %v, ожидается %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && (math.IsNaN(got.Percent) || math.Abs(got.Percent-tt.want) > 1e-10 || got.Graded != tt.graded || got.Criteria != 2) {
				t.Fatalf("итог = %+v, ожидается %.2f, оценено %d/2", got, tt.want, tt.graded)
			}
		})
	}
}

func TestCalculateUnavailable(t *testing.T) {
	t.Parallel()
	// Arrange: курс без рубрикатора.
	for _, tt := range []struct {
		student string
		want    error
	}{
		{"s1", domain.ErrEmptyRubric}, {"unknown", domain.ErrStudentNotFound},
	} {
		t.Run(tt.student, func(t *testing.T) {
			t.Parallel()
			course := &domain.Course{}
			must(t, course.AddStudent("s1", "Анна"))
			// Act: пытаемся получить итог.
			_, err := (domain.Calculator{}).Calculate(course, tt.student)
			// Assert: итог не подменяется произвольным нулём.
			if !errors.Is(err, tt.want) {
				t.Fatalf("ошибка = %v, ожидается %v", err, tt.want)
			}
		})
	}
}

func TestTotalAfterChanges(t *testing.T) {
	t.Parallel()
	// Arrange: полный набор оценок с итогом 85%.
	for _, tt := range []struct {
		name   string
		change func(*domain.Course) error
		want   float64
	}{
		{"исправление", func(c *domain.Course) error { return c.Evaluate("s1", "exam", 10) }, 100},
		{"отмена", func(c *domain.Course) error { return c.CancelEvaluation("s1", "lab") }, 60},
		{"новый критерий", func(c *domain.Course) error { return c.AddCriterion("project", "Проект", 10, 4) }, 42.5},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			course := newCourse(t)
			must(t, course.Evaluate("s1", "exam", 8))
			must(t, course.Evaluate("s1", "lab", 20))
			calculator := domain.Calculator{}
			_, err := calculator.Calculate(course, "s1")
			must(t, err)
			// Act: меняем уже рассчитанный курс.
			must(t, tt.change(course))
			got, err := calculator.Calculate(course, "s1")
			// Assert: расчёт не использует устаревший итог.
			must(t, err)
			if math.Abs(got.Percent-tt.want) > 1e-10 {
				t.Fatalf("итог = %g, ожидается %g", got.Percent, tt.want)
			}
		})
	}
}
