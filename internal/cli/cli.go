// Package cli реализует текстовый интерфейс (REPL) зоопарка.
//
// Пакет зависит только от internal/zoo и стандартной библиотеки.
// Ввод и вывод приходят снаружи (io.Reader/io.Writer) — поэтому REPL
// можно тестировать, подставляя strings.Reader и bytes.Buffer.
package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"example.com/zoo/internal/zoo"
)

const helpText = `Команды:
  add-enclosure <название> <вместимость>    создать вольер
  add-animal <вольер> <имя> <вид> <возраст> добавить животное в вольер
  list                                      показать зоопарк
  help                                      эта справка
  exit                                      выйти`

// Run запускает цикл «команда — ответ» до команды exit или конца ввода.
func Run(z *zoo.Zoo, in io.Reader, out io.Writer) error {
	fmt.Fprintln(out, "Добро пожаловать в зоопарк! Введите help для списка команд.")
	sc := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, "> ")
		if !sc.Scan() { // EOF (Ctrl+D) или ошибка чтения — выходим
			return sc.Err()
		}
		done, err := dispatch(z, sc.Text(), out)
		if err != nil {
			fmt.Fprintln(out, "Ошибка:", err)
		}
		if done {
			fmt.Fprintln(out, "Пока!")
			return nil
		}
	}
}

// dispatch разбирает одну строку и выполняет команду.
// done == true означает «завершить работу».
func dispatch(z *zoo.Zoo, line string, out io.Writer) (done bool, err error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false, nil
	}

	// TODO(ДЗ): новые команды добавляйте сюда — например, remove-animal или feed.
	switch cmd, args := fields[0], fields[1:]; cmd {
	case "exit", "quit":
		return true, nil
	case "help":
		fmt.Fprintln(out, helpText)
	case "add-enclosure":
		return false, addEnclosure(z, args)
	case "add-animal":
		return false, addAnimal(z, args)
	case "list":
		printZoo(z, out)
	default:
		fmt.Fprintf(out, "Неизвестная команда %q. Введите help.\n", cmd)
	}
	return false, nil
}

func addEnclosure(z *zoo.Zoo, args []string) error {
	if len(args) != 2 {
		return errors.New("использование: add-enclosure <название> <вместимость>")
	}
	capacity, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("вместимость должна быть числом: %w", err)
	}
	e, err := zoo.NewEnclosure(args[0], capacity)
	if err != nil {
		return err
	}
	return z.AddEnclosure(e)
}

func addAnimal(z *zoo.Zoo, args []string) error {
	if len(args) != 4 {
		return errors.New("использование: add-animal <вольер> <имя> <вид> <возраст>")
	}
	age, err := strconv.Atoi(args[3])
	if err != nil {
		return fmt.Errorf("возраст должен быть числом: %w", err)
	}
	a, err := zoo.NewAnimal(args[1], args[2], age)
	if err != nil {
		return err
	}
	e, err := z.Enclosure(args[0])
	if err != nil {
		return err
	}
	return e.Add(a)
}

func printZoo(z *zoo.Zoo, out io.Writer) {
	fmt.Fprintf(out, "Зоопарк %q:\n", z.Name())
	for _, e := range z.Enclosures() {
		fmt.Fprintf(out, "  Вольер %q (%d/%d):\n", e.Name(), len(e.Animals()), e.Capacity())
		for _, a := range e.Animals() {
			fmt.Fprintln(out, "   -", a)
		}
	}
}
