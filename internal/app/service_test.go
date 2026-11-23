package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"example.com/zoo/internal/app"
	"example.com/zoo/internal/domain"
	"example.com/zoo/internal/memory"
)

func TestGradeRepository(t *testing.T) {
	t.Parallel()
	// Arrange: порт сообщает ошибки чтения или записи; модель отклоняет неверный балл.
	errStorage := errors.New("хранилище недоступно")
	for _, tt := range []struct {
		name                   string
		score                  float64
		loadErr, saveErr, want error
	}{
		{"успех", 8, nil, nil, nil},
		{"ошибка чтения", 8, errStorage, nil, errStorage},
		{"ошибка записи", 8, nil, errStorage, errStorage},
		{"инвариант", 12, nil, nil, domain.ErrInvalidScore},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			course := &domain.Course{}
			must(t, course.AddStudent("s1", "Анна"))
			must(t, course.AddCriterion("exam", "Экзамен", 10, 1))
			repository := NewMockRepository(t)
			repository.EXPECT().Load().Return(course.Clone(), tt.loadErr).Once()
			if tt.loadErr == nil && tt.score <= 10 {
				repository.EXPECT().Save(mock.MatchedBy(func(saved *domain.Course) bool {
					grade, ok := saved.Evaluation("s1", "exam")
					return ok && grade.Score() == tt.score
				})).Return(tt.saveErr).Once()
			}
			service := newService(repository, NewMockExporter(t))
			// Act: выполняем сценарий через фасад.
			err := service.Grade("s1", "exam", tt.score)
			// Assert: причина сохраняется; после ошибки чтения или модели Save не вызывается.
			if !errors.Is(err, tt.want) {
				t.Fatalf("ошибка = %v, ожидается %v", err, tt.want)
			}
			if _, exists := course.Evaluation("s1", "exam"); exists {
				t.Fatal("изменён сохранённый оригинал")
			}
		})
	}
}

func TestSummary(t *testing.T) {
	t.Parallel()
	// Arrange: группа включает и студента без оценок.
	for _, tt := range []struct {
		name             string
		students, rubric bool
		wantCount        int
		wantErr          error
	}{
		{"пустой курс", false, false, 0, nil},
		{"взвешенная сводка", true, true, 2, nil},
		{"нет критериев", true, false, 0, domain.ErrEmptyRubric},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := newService(&memory.Repository{}, NewMockExporter(t))
			if tt.students {
				must(t, service.AddStudent("s2", "Борис"))
				must(t, service.AddStudent("s1", "Анна"))
			}
			if tt.rubric {
				must(t, service.AddCriterion("exam", "Экзамен", 10, 3))
				must(t, service.AddCriterion("lab", "Практика", 20, 1))
				must(t, service.Grade("s1", "exam", 8))
				must(t, service.Grade("s1", "lab", 20))
			}
			// Act: получаем сводку.
			got, err := service.Summary()
			// Assert: среднее включает пропуски, строки упорядочены.
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ошибка = %v, ожидается %v", err, tt.wantErr)
			}
			if len(got.Results) != tt.wantCount {
				t.Fatalf("строк = %d", len(got.Results))
			}
			if tt.rubric && (got.Average != 42.5 || got.Minimum != 0 || got.Maximum != 85 || got.Results[0].Student.ID() != "s1") {
				t.Fatalf("неверная сводка: %+v", got)
			}
		})
	}
}

func TestExportPort(t *testing.T) {
	t.Parallel()
	// Arrange: ошибки загрузки, расчёта и выгрузки.
	errIO := errors.New("сбой ввода-вывода")
	for _, tt := range []struct {
		name                     string
		loadErr, exportErr, want error
		rubric                   bool
	}{
		{"успех", nil, nil, nil, true},
		{"сбой загрузки", errIO, nil, errIO, true},
		{"нет рубрикатора", nil, nil, domain.ErrEmptyRubric, false},
		{"сбой экспорта", nil, errIO, errIO, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			course := &domain.Course{}
			must(t, course.AddStudent("s1", "Анна"))
			if tt.rubric {
				must(t, course.AddCriterion("exam", "Экзамен", 10, 1))
				must(t, course.Evaluate("s1", "exam", 8))
			}
			repository, exporter := NewMockRepository(t), NewMockExporter(t)
			repository.EXPECT().Load().Return(course, tt.loadErr).Once()
			if tt.loadErr == nil && tt.rubric {
				student, err := course.Student("s1")
				must(t, err)
				exporter.EXPECT().Export("results.csv", []app.Result{{Student: student, Total: domain.Total{Percent: 80, Graded: 1, Criteria: 1}}}).Return(tt.exportErr).Once()
			}
			// Act: экспортируем через фасад.
			err := newService(repository, exporter).Export("results.csv")
			// Assert: выгрузка получает рассчитанные значения, а не весь агрегат.
			if !errors.Is(err, tt.want) {
				t.Fatalf("ошибка = %v, ожидается %v", err, tt.want)
			}
		})
	}
}

func newService(repository app.Repository, exporter app.Exporter) *app.Service {
	return app.New(repository, exporter, domain.NewCalculator(domain.MissingAsZero{}), domain.NewCalculator(domain.RequireAll{}))
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
