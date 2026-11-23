// Приложение ПОС: сборка зависимостей и запуск консоли.
package main

import (
	"fmt"
	"os"

	"example.com/zoo/internal/cli"
)

func main() {
	container := newContainer()
	if err := cli.Run(container.service, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}
