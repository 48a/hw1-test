package cli

import (
	"fmt"
	"strings"
)

func (c console) students(_ []string) (string, error) {
	course, err := c.service.Snapshot()
	if err != nil {
		return "", err
	}
	var lines []string
	for _, student := range course.Students() {
		lines = append(lines, student.ID()+": "+student.Name())
	}
	return listing(lines, "Нет студентов"), nil
}

func (c console) criteria(_ []string) (string, error) {
	course, err := c.service.Snapshot()
	if err != nil {
		return "", err
	}
	var lines []string
	for _, criterion := range course.Criteria() {
		lines = append(lines, fmt.Sprintf("%s: %s; максимум %g; вес %g",
			criterion.ID(), criterion.Name(), criterion.Maximum(), criterion.Weight()))
	}
	return listing(lines, "Нет критериев"), nil
}

func (c console) grades(args []string) (string, error) {
	course, err := c.service.Snapshot()
	if err != nil {
		return "", err
	}
	if _, err := course.Student(args[0]); err != nil {
		return "", err
	}
	var lines []string
	for _, criterion := range course.Criteria() {
		score := "нет оценки"
		if evaluation, exists := course.Evaluation(args[0], criterion.ID()); exists {
			score = fmt.Sprintf("%g/%g", evaluation.Score(), criterion.Maximum())
		}
		lines = append(lines, criterion.ID()+": "+criterion.Name()+"; "+score)
	}
	return listing(lines, "Нет критериев"), nil
}

func (c console) summary(_ []string) (string, error) {
	summary, err := c.service.Summary()
	if err != nil {
		return "", err
	}
	if len(summary.Results) == 0 {
		return "Нет студентов", nil
	}
	var lines []string
	for _, result := range summary.Results {
		lines = append(lines, fmt.Sprintf("%s: %s; %.2f%%; оценено %d/%d",
			result.Student.ID(), result.Student.Name(), result.Total.Percent, result.Total.Graded, result.Total.Criteria))
	}
	lines = append(lines, fmt.Sprintf("Среднее: %.2f%%; минимум: %.2f%%; максимум: %.2f%%",
		summary.Average, summary.Minimum, summary.Maximum))
	return strings.Join(lines, "\n"), nil
}

func listing(lines []string, empty string) string {
	if len(lines) == 0 {
		return empty
	}
	return strings.Join(lines, "\n")
}
