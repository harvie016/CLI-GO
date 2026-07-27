package main

import (
	"fmt"
	"os"

	"mws/internal/cli"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		if cli.IsHelpShown(err) {
			return
		}
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

// run разбирает верхний уровень команд. Ошибку возвращаем наверх, чтобы весь
// вывод в stderr и код возврата были в одном месте.
func run(args []string) error {
	// Запуск без аргументов — показываем справку, а не молчим.
	if len(args) == 0 {
		cli.Help()
		return nil
	}

	switch args[0] {
	case "profile":
		return cli.Profile(args[1:])
	case "help", "-h", "--help":
		cli.Help()
		return nil
	default:
		return fmt.Errorf("неизвестная команда %q; список команд: mws help", args[0])
	}
}
