// Точка входа приложения «Зоопарк»: сборка зависимостей и запуск REPL.
package main

import (
	"fmt"
	"os"

	"example.com/zoo/internal/cli"
	"example.com/zoo/internal/zoo"
)

func main() {
	z := zoo.New("Домашний зоопарк")
	if err := cli.Run(z, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
}
