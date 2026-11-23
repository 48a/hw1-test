// Package csvexport выгружает итоги в CSV-файл.
package csvexport

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"example.com/zoo/internal/app"
)

const (
	decimalPlaces = 2
	fileMode      = 0o600
)

// Exporter — файловый адаптер порта выгрузки итогов.
type Exporter struct{}

// Export записывает UTF-8 CSV в новый файл, не заменяя существующий.
func (Exporter) Export(path string, results []app.Result) error {
	root, err := os.OpenRoot(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("открытие каталога: %w", err)
	}
	defer root.Close()
	file, err := root.OpenFile(filepath.Base(path), os.O_WRONLY|os.O_CREATE|os.O_EXCL, fileMode)
	if err != nil {
		return fmt.Errorf("создание файла: %w", err)
	}
	writer := csv.NewWriter(file)
	records := [][]string{{"student_id", "name", "total_percent", "graded", "criteria"}}
	for _, result := range results {
		records = append(records, []string{
			result.Student.ID(), result.Student.Name(),
			strconv.FormatFloat(result.Total.Percent, 'f', decimalPlaces, 64),
			strconv.Itoa(result.Total.Graded), strconv.Itoa(result.Total.Criteria),
		})
	}
	err = errors.Join(writer.WriteAll(records), file.Close())
	if err != nil {
		return fmt.Errorf("запись CSV: %w", err)
	}
	return nil
}
