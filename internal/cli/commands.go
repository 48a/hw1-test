package cli

import (
	"fmt"
	"strconv"
	"strings"
)

func (c console) student(args []string) (string, error) {
	return "Студент добавлен", c.service.AddStudent(args[0], strings.Join(args[1:], " "))
}

func (c console) criterion(args []string) (string, error) {
	maximum, err := number(args[1])
	if err != nil {
		return "", err
	}
	weight, err := number(args[2])
	if err != nil {
		return "", err
	}
	const nameStart = 3
	return "Критерий добавлен", c.service.AddCriterion(args[0], strings.Join(args[nameStart:], " "), maximum, weight)
}

func (c console) grade(args []string) (string, error) {
	score, err := number(args[2])
	if err != nil {
		return "", err
	}
	return "Оценка сохранена", c.service.Grade(args[0], args[1], score)
}

func (c console) ungrade(args []string) (string, error) {
	return "Оценка отменена", c.service.Ungrade(args[0], args[1])
}

func (c console) export(args []string) (string, error) {
	return "Экспорт завершён", c.service.Export(strings.Join(args, " "))
}

func (c console) total(args []string) (string, error) {
	final := len(args) > 1
	if final && args[1] != "final" {
		return "", fmt.Errorf("%w: total <студент> [final]", errUsage)
	}
	total, err := c.service.Total(args[0], final)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s: %.2f%%; оценено %d/%d", args[0], total.Percent, total.Graded, total.Criteria), nil
}

func number(raw string) (float64, error) {
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("введите число с точкой: %w", err)
	}
	return value, nil
}
