package cli_test

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/iotest"

	"example.com/zoo/internal/app"
	"example.com/zoo/internal/cli"
	"example.com/zoo/internal/csvexport"
	"example.com/zoo/internal/domain"
	"example.com/zoo/internal/memory"
)

const setup = "student s1 Анна Иванова\ncriterion exam 10 3 Экзамен\n"

func TestRunWorkflow(t *testing.T) {
	t.Parallel()
	// Arrange: сквозные сценарии с реальными адаптерами памяти и CSV.
	for _, tt := range []struct{ name, edits, percent, graded string }{
		{"полный итог", "", "85.00", "2"},
		{"исправление и отмена", "grade s1 exam 9\nungrade s1 lab\n", "67.50", "1"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "итоги курса.csv")
			input := setup + "student s2 Борис Петров\ncriterion lab 20 1 Практика\n" +
				"grades s1\ngrade s1 exam 12\ngrade s1 exam 8\ntotal s1 final\n" +
				"grade s1 lab 20\ntotal s1 final\n" + tt.edits +
				"total s1\nsummary\nexport " + path + "\nexit\n"
			var output bytes.Buffer
			// Act: выполняем консольный сеанс.
			must(t, cli.Run(newService(), strings.NewReader(input), &output))
			data, err := os.ReadFile(path)
			must(t, err)
			records, err := csv.NewReader(strings.NewReader(string(data))).ReadAll()
			// Assert: ошибки не прервали сеанс; экспорт содержит актуальный итог.
			must(t, err)
			for _, fragment := range []string{"нет оценки", "максимум 10", "выставлены не все оценки", "s1: 85.00%; оценено 2/2", "Экспорт завершён", "До свидания!"} {
				if !strings.Contains(output.String(), fragment) {
					t.Errorf("нет %q в выводе:\n%s", fragment, output.String())
				}
			}
			if len(records) != 3 || !slices.Equal(records[1], []string{"s1", "Анна Иванова", tt.percent, tt.graded, "2"}) || !slices.Equal(records[2], []string{"s2", "Борис Петров", "0.00", "0", "2"}) {
				t.Fatalf("неверный CSV: %v", records)
			}
		})
	}
}

func TestRunCommands(t *testing.T) {
	t.Parallel()
	// Arrange: команды просмотра, ошибки формата и инвариантов.
	for _, tt := range []struct{ command, want string }{
		{"help", "Команды:"},
		{"students", "s1: Анна Иванова"},
		{"criteria", "exam: Экзамен; максимум 10; вес 3"},
		{"grades s1", "exam: Экзамен; нет оценки"},
		{"total s1", "s1: 0.00%; оценено 0/1"},
		{"summary", "Среднее: 0.00%"},
		{"student", "student <ID> <имя>"},
		{"students extra", "неверный формат команды"},
		{"criterion bad nope 1 Название", "введите число с точкой"},
		{"criterion bad 10 nope Название", "введите число с точкой"},
		{"grade s1 exam nope", "введите число с точкой"},
		{"grade s1 exam NaN", "балл должен быть конечным"},
		{"grade s1 exam +Inf", "балл должен быть конечным"},
		{"grade s1 missing 1", "критерий не найден"},
		{"grade missing exam 1", "студент не найден"},
		{"grades missing", "студент не найден"},
		{"ungrade s1 exam", "оценка не выставлена"},
		{"student s1 Другое имя", "ID уже существует"},
		{"total s1 wrong", "total <студент> [final]"},
		{"unknown", "неизвестная команда"},
		{"exit extra", "неверный формат команды"},
	} {
		t.Run(tt.command, func(t *testing.T) {
			t.Parallel()
			var output bytes.Buffer
			// Act: после проверяемой команды запрашиваем список и выходим.
			err := cli.Run(newService(), strings.NewReader(setup+tt.command+"\nstudents\nexit\n"), &output)
			// Assert: есть нужный ответ и продолжение сеанса.
			must(t, err)
			if !strings.Contains(output.String(), tt.want) || !strings.HasSuffix(output.String(), "s1: Анна Иванова\n> До свидания!\n") {
				t.Fatalf("неверный вывод:\n%s", output.String())
			}
		})
	}
}

func TestRunEmpty(t *testing.T) {
	t.Parallel()
	// Arrange: пустой курс и завершение через EOF.
	for _, tt := range []struct{ input, want string }{
		{"\n", "ПОС."},
		{"students\n", "Нет студентов"},
		{"criteria\n", "Нет критериев"},
		{"summary\n", "Нет студентов"},
		{"student s1 Анна\ngrades s1\n", "Нет критериев"},
		{"student s1 Анна\ntotal s1\n", "рубрикатор пуст"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			var output bytes.Buffer
			// Act: читаем команды до EOF.
			err := cli.Run(newService(), strings.NewReader(tt.input), &output)
			// Assert: EOF не считается ошибкой.
			must(t, err)
			if !strings.Contains(output.String(), tt.want) {
				t.Fatalf("неверный вывод: %s", output.String())
			}
		})
	}
}

func TestRunIOErrors(t *testing.T) {
	t.Parallel()
	// Arrange: стандартные источники ошибок чтения и записи.
	for _, operation := range []string{"чтение", "запись"} {
		t.Run(operation, func(t *testing.T) {
			t.Parallel()
			want := errors.New("сбой потока")
			input := iotest.ErrReader(want)
			output := io.Discard
			if operation == "запись" {
				reader, writer := io.Pipe()
				must(t, reader.CloseWithError(want))
				t.Cleanup(func() { must(t, writer.Close()) })
				input, output = strings.NewReader("exit\n"), writer
			}
			// Act: запускаем REPL с недоступным потоком.
			err := cli.Run(newService(), input, output)
			// Assert: ошибка доходит до точки входа.
			if !errors.Is(err, want) {
				t.Fatalf("ошибка = %v, ожидается %v", err, want)
			}
		})
	}
}

func newService() *app.Service {
	return app.New(&memory.Repository{}, csvexport.Exporter{}, domain.NewCalculator(domain.MissingAsZero{}), domain.NewCalculator(domain.RequireAll{}))
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
