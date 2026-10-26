package cli_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"example.com/zoo/internal/cli"
	"example.com/zoo/internal/domain"
)

func newPark(t *testing.T) *domain.Park {
	t.Helper()
	park, err := domain.NewPark(domain.BasicServiceCenter{})
	if err != nil {
		t.Fatal(err)
	}
	return park
}

func TestRun(t *testing.T) {
	t.Parallel()
	// Arrange: успешные команды и ошибки пользовательского ввода.
	tests := []struct {
		name, input, want string
		total             int
	}{
		{"выход", "0\n1\n", "До свидания!", 0},
		{"пустой отчёт", "3\n0\n", "Транспорт: 0; вещи: 0; всего: 0; Суточная энергия всего баланса: 0.00 кВт·ч", 0},
		{"пустая подборка", "4\n0\n", "Список пуст.", 0},
		{"пустой инвентарь", "5\n0\n", "Весь инвентарь:\nСписок пуст.", 0},
		{"неизвестная команда", "9\n0\n", "Неизвестная команда. Введите число от 0 до 5.", 0},
		{"текст вместо команды", "abc\n0\n", "Неизвестная команда.", 0},
		{"неизвестный транспорт", "1\n4\n0\n", "Не добавлено: неизвестный вид: 4", 0},
		{"неизвестная вещь", "2\n4\n0\n", "Не добавлено: неизвестный вид: 4", 0},
		{"пустое название", "2\n1\n  \n0\n", "Не добавлено: название обязательно", 0},
		{"пустой номер", "2\n1\nШлем\n\n0\n", "Не добавлено: название обязательно", 0},
		{"нечисловой номер", "2\n1\nШлем\nabc\n0\n", "инвентарный номер должен быть положительным целым числом", 0},
		{"нулевой номер", "2\n1\nШлем\n0\n0\n", "инвентарный номер должен быть положительным целым числом", 0},
		{"простота текстом", "1\n2\nВелосипед\n1\nabc\n0\n", "Не добавлено: простота должна быть целым числом от 1 до 10", 0},
		{"простота дробью", "1\n2\nВелосипед\n1\n6.5\n0\n", "Не добавлено: простота должна быть целым числом от 1 до 10", 0},
		{"простота вне границ", "1\n2\nВелосипед\n1\n11\n0\n", "Не добавлено: простота должна быть целым числом от 1 до 10", 0},
		{"неверное состояние", "1\n2\nВелосипед\n1\n6\n3\n0\n", "Не добавлено: техническое состояние должно быть 1", 0},
		{"отказ техосмотра", "1\n1\nСамокат\n1\n9\n2\n1.5\n0\n", "Не добавлено: техосмотр отклонил транспорт", 0},
		{"расход текстом", "2\n3\nШкаф\n1\nabc\n0\n", "Не добавлено: суточный расход энергии", 0},
		{"расход с запятой", "2\n3\nШкаф\n1\n1,5\n0\n", "Не добавлено: суточный расход энергии", 0},
		{"отрицательный расход", "1\n1\nСамокат\n1\n6\n1\n-1\n0\n", "Не добавлено: суточный расход энергии", 0},
		{"NaN", "2\n3\nШкаф\n1\nNaN\n0\n", "Не добавлено: суточный расход энергии", 0},
		{"бесконечность", "2\n3\nШкаф\n1\nInf\n0\n", "Не добавлено: суточный расход энергии", 0},
		{"нормализация", "2\n1\n  Шлем с пробелами  \n001\n0\n", "Принято: Шлем с пробелами [1]", 1},
		{"повтор номера", "2\n1\nШлем\n1\n2\n2\nДок\n001\n0\n", "Не добавлено: инвентарный номер уже занят", 1},
		{"восстановление после ошибки", "1\n4\n2\n1\nШлем\n1\n5\n0\n", "Весь инвентарь:\n- Шлем [1]", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			park := newPark(t)
			var output bytes.Buffer
			// Act: выполняем ввод как в консоли.
			err := cli.Run(park, strings.NewReader(tt.input), &output)
			// Assert: ошибка ввода показана пользователю, программа дошла до выхода.
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(output.String(), tt.want) || !strings.Contains(output.String(), "До свидания!") {
				t.Errorf("вывод:\n%s\nожидался фрагмент %q и завершение", &output, tt.want)
			}
			if park.Report().Total != tt.total {
				t.Errorf("баланс = %+v, ожидалось объектов %d", park.Report(), tt.total)
			}
		})
	}
}

func TestDemo(t *testing.T) {
	t.Parallel()
	// Arrange: демонстрация из README содержит все виды и основные ошибки.
	input, err := os.ReadFile("../../docs/demo.txt")
	if err != nil {
		t.Fatal(err)
	}
	park := newPark(t)
	var output bytes.Buffer
	// Act: выполняем тот же сценарий, что и make demo.
	err = cli.Run(park, bytes.NewReader(input), &output)
	// Assert: приняты шесть объектов, ошибки не меняют баланс.
	if err != nil {
		t.Fatal(err)
	}
	if want := (domain.Report{Transport: 3, Items: 3, Total: 6, Energy: 6.75}); park.Report() != want {
		t.Fatalf("отчёт = %+v, ожидался %+v", park.Report(), want)
	}
	for _, fragment := range []string{
		"Транспорт: 3; вещи: 3; всего: 6; Суточная энергия всего баланса: 6.75 кВт·ч",
		"Техника для новичков (простота >= 6):\n- City Pro [101], простота: 6\n- Электровелосипед [103], простота: 10",
		"Весь инвентарь:\n- City Pro [101]\n- Городской велосипед [102]\n- Электровелосипед [103]\n- Шлем [201]\n- Док-станция без зарядки [202]\n- Зарядный шкаф [203]",
		"Не добавлено: техосмотр отклонил транспорт",
		"Не добавлено: инвентарный номер уже занят",
		"Не добавлено: неизвестный вид: 4",
		"Не добавлено: простота должна быть целым числом от 1 до 10",
		"Не добавлено: техническое состояние должно быть 1",
		"Не добавлено: суточный расход энергии",
		"Неизвестная команда. Введите число от 0 до 5.",
		"До свидания!",
	} {
		if !strings.Contains(output.String(), fragment) {
			t.Errorf("в выводе нет %q", fragment)
		}
	}
	if strings.Count(output.String(), "Принято:") != 6 {
		t.Errorf("лишняя или пропущенная приёмка:\n%s", &output)
	}
}

func TestIncompleteInput(t *testing.T) {
	t.Parallel()
	// Arrange: обрывы ввода перед каждым полем приёмки транспорта и вещей.
	for _, scenario := range []string{"1\n1\nСамокат\n1\n6\n1\n1.5", "2\n3\nШкаф\n2\n3"} {
		fields := strings.Split(scenario, "\n")
		for i := range len(fields) {
			park := newPark(t)
			input := strings.Join(fields[:i], "\n")
			if i > 0 {
				input += "\n"
			}
			// Act: завершаем поток посреди команды.
			err := cli.Run(park, strings.NewReader(input), io.Discard)
			// Assert: EOF штатно завершает приложение, частичной приёмки нет.
			if err != nil || park.Report().Total != 0 {
				t.Errorf("ввод %q: ошибка = %v, отчёт = %+v", input, err, park.Report())
			}
		}
	}
}

var errBrokenStream = errors.New("поток недоступен")

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) { return 0, errBrokenStream }

type brokenWriter struct{ writesLeft int }

func (w *brokenWriter) Write(data []byte) (int, error) {
	if w.writesLeft == 0 {
		return 0, errBrokenStream
	}
	w.writesLeft--
	return len(data), nil
}

func TestStreamErrors(t *testing.T) {
	t.Parallel()
	// Arrange: ошибки чтения и вывода на разных шагах диалога.
	tests := []struct {
		name string
		in   func() io.Reader
		out  func() io.Writer
	}{
		{"чтение команды", func() io.Reader { return brokenReader{} }, func() io.Writer { return io.Discard }},
		{"чтение поля", func() io.Reader { return io.MultiReader(strings.NewReader("1\n1\n"), brokenReader{}) }, func() io.Writer { return io.Discard }},
		{"вывод меню", func() io.Reader { return strings.NewReader("0\n") }, func() io.Writer { return &brokenWriter{} }},
		{"вывод поля", func() io.Reader { return strings.NewReader("1\n1\n") }, func() io.Writer { return &brokenWriter{writesLeft: 1} }},
		{"финальный вывод", func() io.Reader { return strings.NewReader("0\n") }, func() io.Writer { return &brokenWriter{writesLeft: 1} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			park := newPark(t)
			// Act: запускаем консоль с повреждённым потоком.
			err := cli.Run(park, tt.in(), tt.out())
			// Assert: ошибка не потеряна и доступна через errors.Is.
			if !errors.Is(err, errBrokenStream) || park.Report().Total != 0 {
				t.Fatalf("ошибка = %v, отчёт = %+v", err, park.Report())
			}
		})
	}
}
