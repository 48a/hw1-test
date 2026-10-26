// Приложение «ВышКат»: сборка зависимостей и запуск консоли.
package main

import (
	"fmt"
	"os"

	"example.com/zoo/internal/cli"
	"example.com/zoo/internal/domain"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

func run() error {
	park, err := domain.NewPark(domain.BasicServiceCenter{})
	if err != nil {
		return err
	}
	return cli.Run(park, os.Stdin, os.Stdout)
}
