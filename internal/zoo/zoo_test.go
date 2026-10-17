package zoo_test

import (
	"errors"
	"testing"

	"example.com/zoo/internal/zoo"
)

// TestNewAnimal — образец table-driven теста:
// таблица случаев + один цикл запуска. Добавляйте строки — не новые функции.
func TestNewAnimal(t *testing.T) {
	t.Parallel()

	// Arrange: таблица случаев.
	tests := []struct {
		name    string // название случая, увидите его в go test -run
		animal  string // имя животного
		species string
		age     int
		wantErr error // ожидаемая ошибка (сравниваем через errors.Is)
	}{
		{name: "корректное животное", animal: "Буся", species: "капибара", age: 3, wantErr: nil},
		{name: "без имени", animal: "", species: "капибара", age: 3, wantErr: zoo.ErrInvalidAnimal},
		{name: "без вида", animal: "Буся", species: "", age: 3, wantErr: zoo.ErrInvalidAnimal},
		{name: "отрицательный возраст", animal: "Буся", species: "капибара", age: -1, wantErr: zoo.ErrInvalidAnimal},
		{name: "нулевой возраст — ок", animal: "Буся", species: "капибара", age: 0, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act: создаём животное.
			got, err := zoo.NewAnimal(tt.animal, tt.species, tt.age)

			// Assert: ошибка и поля совпадают с ожиданием.
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewAnimal() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got.Name != tt.animal {
				t.Errorf("NewAnimal() name = %q, want %q", got.Name, tt.animal)
			}
		})
	}
}

func TestEnclosureAddRemove(t *testing.T) {
	t.Parallel()

	// Arrange: вольер на одно место и два животных.
	e, err := zoo.NewEnclosure("Тропики", 1)
	if err != nil {
		t.Fatalf("NewEnclosure() error = %v", err)
	}

	capybara, err := zoo.NewAnimal("Буся", "капибара", 3)
	if err != nil {
		t.Fatalf("NewAnimal() error = %v", err)
	}
	otter, err := zoo.NewAnimal("Шнырь", "выдра", 2)
	if err != nil {
		t.Fatalf("NewAnimal() error = %v", err)
	}

	// Act + Assert: первая капибара помещается, для выдры места нет.
	if err := e.Add(capybara); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if err := e.Add(otter); !errors.Is(err, zoo.ErrEnclosureFull) {
		t.Fatalf("Add() в полный вольер error = %v, want %v", err, zoo.ErrEnclosureFull)
	}

	// Act + Assert: удаление возвращает животное; повторное удаление — ошибку.
	got, err := e.Remove("Буся")
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if got != capybara {
		t.Errorf("Remove() = %v, want %v", got, capybara)
	}
	if _, err := e.Remove("Буся"); !errors.Is(err, zoo.ErrAnimalNotFound) {
		t.Fatalf("Remove() ещё раз error = %v, want %v", err, zoo.ErrAnimalNotFound)
	}
}

func TestZooEnclosuresSorted(t *testing.T) {
	t.Parallel()

	// Arrange: зоопарк с вольерами, добавленными не по алфавиту.
	z := zoo.New("Тестовый")
	for _, name := range []string{"Саванна", "Аквариум", "Тропики"} {
		e, err := zoo.NewEnclosure(name, 2)
		if err != nil {
			t.Fatalf("NewEnclosure(%q) error = %v", name, err)
		}
		if err := z.AddEnclosure(e); err != nil {
			t.Fatalf("AddEnclosure(%q) error = %v", name, err)
		}
	}

	// Act: запрашиваем список вольеров.
	enclosures := z.Enclosures()

	// Assert: названия отсортированы по алфавиту.
	var names []string
	for _, e := range enclosures {
		names = append(names, e.Name())
	}
	want := []string{"Аквариум", "Саванна", "Тропики"}
	if len(names) != len(want) {
		t.Fatalf("Enclosures() вернул %d вольеров, want %d", len(names), len(want))
	}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("Enclosures()[%d] = %q, want %q", i, names[i], want[i])
		}
	}
}

// TestFeedAll — тест с моком. Мок MockNotifier сгенерирован mockery
// (make mocks, конфиг .mockery.yml) из интерфейса zoo.Notifier:
// настоящий канал уведомлений в тестах не нужен, а поведение
// (что, сколько раз и с какими аргументами вызвано) — проверяем.
func TestFeedAll(t *testing.T) {
	t.Parallel()

	// Arrange (общий): два животных.
	capybara, err := zoo.NewAnimal("Буся", "капибара", 3)
	if err != nil {
		t.Fatalf("NewAnimal() error = %v", err)
	}
	otter, err := zoo.NewAnimal("Шнырь", "выдра", 2)
	if err != nil {
		t.Fatalf("NewAnimal() error = %v", err)
	}
	animals := []zoo.Animal{capybara, otter}

	t.Run("уведомление отправляется про каждое животное", func(t *testing.T) {
		t.Parallel()

		// Arrange: мок ждёт ровно два вызова с конкретными сообщениями.
		// Ожидания проверяются автоматически в t.Cleanup() при конце теста.
		n := NewMockNotifier(t)
		n.EXPECT().Notify("кормлю: Буся").Return(nil).Once()
		n.EXPECT().Notify("кормлю: Шнырь").Return(nil).Once()

		// Act: кормим обоих.
		err := zoo.FeedAll(n, animals)
		// Assert: ошибки нет; вызовы мока проверит t.Cleanup().
		if err != nil {
			t.Fatalf("FeedAll() error = %v", err)
		}
	})

	t.Run("ошибка уведомления прерывает кормление", func(t *testing.T) {
		t.Parallel()

		// Arrange: мок «падает» на первом же вызове.
		errDown := errors.New("канал уведомлений недоступен")
		n := NewMockNotifier(t)
		n.EXPECT().Notify("кормлю: Буся").Return(errDown).Once()

		// Act: кормим обоих, но второй вызов не должен случиться.
		err := zoo.FeedAll(n, animals)

		// Assert: ошибка проброшена наружу (обёрнута через %w).
		// То, что Notify вызвали ровно один раз, подтвердит t.Cleanup() мока.
		if !errors.Is(err, errDown) {
			t.Fatalf("FeedAll() error = %v, want %v", err, errDown)
		}
	})
}
