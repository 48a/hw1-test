// Package cli реализует текстовый интерфейс ПОС.
package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"example.com/zoo/internal/app"
	"example.com/zoo/internal/domain"
)

const helpText = `Команды:
  student <ID> <имя>                  добавить студента
  students                           список студентов
  criterion <ID> <максимум> <вес> <название>
  criteria                           список критериев
  grade <студент> <критерий> <балл>    выставить или исправить
  grades <студент>                    посмотреть оценки
  ungrade <студент> <критерий>         отменить оценку
  total <студент> [final]             текущий или окончательный итог
  summary                            итоги и статистика по группе
  export <путь.csv>                   записать новый CSV-файл
  help                               справка
  exit                               выход`

var errUsage = errors.New("неверный формат команды")

const unlimited = -1

type command struct {
	usage            string
	minimum, maximum int
	run              func([]string) (string, error)
}

type console struct{ service *app.Service }

// Run запускает REPL до exit или EOF; ошибки команд не прерывают сеанс.
func Run(service *app.Service, in io.Reader, out io.Writer) error {
	if _, err := fmt.Fprintln(out, "ПОС. Данные хранятся до выхода. Введите help."); err != nil {
		return fmt.Errorf("вывод: %w", err)
	}
	commands := console{service: service}.commands()
	scanner := bufio.NewScanner(in)
	for {
		if _, err := fmt.Fprint(out, "> "); err != nil {
			return fmt.Errorf("вывод: %w", err)
		}
		if !scanner.Scan() {
			return scanner.Err()
		}
		response, done, err := execute(commands, scanner.Text())
		if err != nil {
			response = fmt.Sprintf("Ошибка: %v. %s", err, hint(err))
		}
		if response != "" {
			if _, err := fmt.Fprintln(out, response); err != nil {
				return fmt.Errorf("вывод: %w", err)
			}
		}
		if done {
			return nil
		}
	}
}

func (c console) commands() map[string]command {
	return map[string]command{
		"student":   {"student <ID> <имя>", 2, unlimited, c.student},
		"students":  {"students", 0, 0, c.students},
		"criterion": {"criterion <ID> <максимум> <вес> <название>", 4, unlimited, c.criterion},
		"criteria":  {"criteria", 0, 0, c.criteria},
		"grade":     {"grade <студент> <критерий> <балл>", 3, 3, c.grade},
		"grades":    {"grades <студент>", 1, 1, c.grades},
		"ungrade":   {"ungrade <студент> <критерий>", 2, 2, c.ungrade},
		"total":     {"total <студент> [final]", 1, 2, c.total},
		"summary":   {"summary", 0, 0, c.summary},
		"export":    {"export <путь.csv>", 1, unlimited, c.export},
		"help":      {"help", 0, 0, func(_ []string) (string, error) { return helpText, nil }},
	}
}

func execute(commands map[string]command, line string) (string, bool, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "", false, nil
	}
	if fields[0] == "exit" && len(fields) == 1 {
		return "До свидания!", true, nil
	}
	cmd, exists := commands[fields[0]]
	if !exists {
		return "", false, fmt.Errorf("%w: неизвестная команда %q", errUsage, fields[0])
	}
	args := fields[1:]
	if len(args) < cmd.minimum || (cmd.maximum != unlimited && len(args) > cmd.maximum) {
		return "", false, fmt.Errorf("%w: %s", errUsage, cmd.usage)
	}
	response, err := cmd.run(args)
	return response, false, err
}

func hint(err error) string {
	switch {
	case errors.Is(err, domain.ErrStudentNotFound):
		return "Проверьте ID командой students."
	case errors.Is(err, domain.ErrCriterionNotFound), errors.Is(err, domain.ErrInvalidScore):
		return "Проверьте критерии командой criteria."
	case errors.Is(err, domain.ErrEmptyRubric):
		return "Добавьте критерий командой criterion."
	case errors.Is(err, domain.ErrIncomplete), errors.Is(err, domain.ErrEvaluationNotFound):
		return "Посмотрите grades и выставьте оценки командой grade."
	case errors.Is(err, os.ErrExist):
		return "Укажите новое имя файла."
	default:
		return "Проверьте данные; help — формат команд."
	}
}
