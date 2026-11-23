package csvexport_test

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"example.com/zoo/internal/app"
	"example.com/zoo/internal/csvexport"
	"example.com/zoo/internal/domain"
)

func TestExportCSV(t *testing.T) {
	t.Parallel()
	// Arrange: пустой отчёт и имя с символами, требующими CSV-экранирования.
	for _, tt := range []struct {
		name, studentName string
		empty             bool
	}{
		{"обычный", "Анна Иванова", false},
		{"кавычки и запятые", "Анна, \"Иванова\"", false},
		{"пустой", "", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "итоги курса.csv")
			var results []app.Result
			want := [][]string{{"student_id", "name", "total_percent", "graded", "criteria"}}
			if !tt.empty {
				course := &domain.Course{}
				must(t, course.AddStudent("s1", tt.studentName))
				student, err := course.Student("s1")
				must(t, err)
				results = []app.Result{{Student: student, Total: domain.Total{Percent: 85.12345, Graded: 2, Criteria: 3}}}
				want = append(want, []string{"s1", tt.studentName, "85.12", "2", "3"})
			}
			// Act: выгружаем файл и читаем стандартным CSV-парсером.
			must(t, (csvexport.Exporter{}).Export(path, results))
			data, err := os.ReadFile(path)
			must(t, err)
			got, err := csv.NewReader(strings.NewReader(string(data))).ReadAll()
			// Assert: точные поля, экранирование и округление.
			must(t, err)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("CSV = %v, ожидается %v", got, want)
			}
		})
	}
}

func TestExportErrors(t *testing.T) {
	t.Parallel()
	// Arrange: существующий файл и отсутствующий каталог.
	for _, tt := range []struct {
		name   string
		exists bool
		want   error
	}{
		{"файл существует", true, os.ErrExist},
		{"каталог отсутствует", false, os.ErrNotExist},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "results.csv")
			if tt.exists {
				must(t, os.WriteFile(path, []byte("старые данные"), 0o600))
			} else {
				path = filepath.Join(path, "results.csv")
			}
			// Act: пытаемся выгрузить отчёт.
			err := (csvexport.Exporter{}).Export(path, nil)
			// Assert: причина доступна через errors.Is; прежний файл не затронут.
			if !errors.Is(err, tt.want) {
				t.Fatalf("ошибка = %v, ожидается %v", err, tt.want)
			}
			if tt.exists {
				data, err := os.ReadFile(path)
				must(t, err)
				if string(data) != "старые данные" {
					t.Fatal("существующий файл изменён")
				}
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
