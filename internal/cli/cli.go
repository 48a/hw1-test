// Package cli реализует пошаговое текстовое меню приложения «ВышКат».
package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"example.com/zoo/internal/domain"
)

const menu = `ВышКат — учёт прокатного парка
1 — Добавить транспорт
2 — Добавить вещь
3 — Отчёт
4 — Техника для новичков
5 — Весь инвентарь
0 — Выход`

type console struct {
	scanner *bufio.Scanner
	out     *bufio.Writer
	park    *domain.Park
}

// Run выполняет команды до выхода или конца ввода.
func Run(park *domain.Park, in io.Reader, out io.Writer) (err error) {
	c := console{scanner: bufio.NewScanner(in), out: bufio.NewWriter(out), park: park}
	defer func() { err = errors.Join(err, c.out.Flush()) }()
	fmt.Fprintln(c.out, menu)
	err = c.loop()
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func (c *console) loop() error {
	for {
		command, err := c.read("Команда (0–5): ")
		if err != nil {
			return err
		}
		if command == "0" {
			fmt.Fprintln(c.out, "До свидания!")
			return nil
		}
		if err := c.execute(command); err != nil {
			if c.scanner.Err() != nil || errors.Is(err, io.EOF) {
				return err
			}
			if flushErr := c.out.Flush(); flushErr != nil {
				return errors.Join(err, flushErr)
			}
			fmt.Fprintf(c.out, "Не добавлено: %v. Повторите команду 1 или 2 с корректными данными.\n", err)
		}
	}
}

func (c *console) read(prompt string) (string, error) {
	fmt.Fprintln(c.out, strings.TrimSpace(prompt))
	if err := c.out.Flush(); err != nil {
		return "", err
	}
	if !c.scanner.Scan() {
		if err := c.scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(c.scanner.Text()), nil
}

func (c *console) execute(command string) error {
	switch command {
	case "1":
		return c.addTransport()
	case "2":
		return c.addItem()
	case "3":
		r := c.park.Report()
		fmt.Fprintf(c.out, "Транспорт: %d; вещи: %d; всего: %d; Суточная энергия всего баланса: %.2f кВт·ч\n",
			r.Transport, r.Items, r.Total, r.Energy)
	case "4":
		c.printBeginners()
	case "5":
		c.printInventory()
	default:
		fmt.Fprintln(c.out, "Неизвестная команда. Введите число от 0 до 5.")
	}
	return nil
}

func (c *console) printBeginners() {
	fmt.Fprintln(c.out, "Техника для новичков (простота >= 6):")
	transport := c.park.Beginners()
	if len(transport) == 0 {
		fmt.Fprintln(c.out, "Список пуст.")
	}
	for _, item := range transport {
		fmt.Fprintf(c.out, "- %s [%s], простота: %d\n", item.Name(), item.InventoryNumber(), item.BeginnerEase())
	}
}

func (c *console) printInventory() {
	fmt.Fprintln(c.out, "Весь инвентарь:")
	items := c.park.Inventory()
	if len(items) == 0 {
		fmt.Fprintln(c.out, "Список пуст.")
	}
	for _, item := range items {
		fmt.Fprintf(c.out, "- %s [%s]\n", item.Name(), item.InventoryNumber())
	}
}
