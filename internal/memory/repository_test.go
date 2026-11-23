package memory_test

import (
	"testing"

	"example.com/zoo/internal/domain"
	"example.com/zoo/internal/memory"
)

func TestRepositoryIsolation(t *testing.T) {
	t.Parallel()
	// Arrange: проверяем обе границы копирования и списки агрегата.
	for _, operation := range []string{"после Save", "после Load", "списки"} {
		t.Run(operation, func(t *testing.T) {
			t.Parallel()
			repository := &memory.Repository{}
			course, err := repository.Load()
			must(t, err)
			must(t, course.AddStudent("s2", "Борис"))
			must(t, course.AddStudent("s1", "Анна"))
			must(t, course.AddCriterion("exam", "Экзамен", 10, 1))
			must(t, course.Evaluate("s1", "exam", 8))
			must(t, repository.Save(course))
			// Act: изменяем исходный объект, снимок или возвращённые списки.
			if operation == "после Load" {
				course, err = repository.Load()
				must(t, err)
			}
			if operation == "списки" {
				course.Students()[0] = domain.Student{}
				course.Criteria()[0] = domain.Criterion{}
			} else {
				must(t, course.Evaluate("s1", "exam", 0))
				must(t, course.AddStudent("s3", "Вера"))
				must(t, course.AddCriterion("lab", "Практика", 20, 2))
			}
			// Assert: сохранённый курс остаётся прежним, списки упорядочены.
			stored, err := repository.Load()
			must(t, err)
			grade, ok := stored.Evaluation("s1", "exam")
			students, criteria := stored.Students(), stored.Criteria()
			if !ok || grade.Score() != 8 || len(students) != 2 || len(criteria) != 1 || students[0].ID() != "s1" || criteria[0].Maximum() != 10 {
				t.Fatalf("сохранённые данные изменились: %v, %v, %v", students, criteria, grade)
			}
		})
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
