// Package cli разбирает аргументы командной строки и вызывает слой хранения.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
)

// Profile — маршрутизация внутри группы команд «mws profile ...».
func Profile(args []string) error {
	if len(args) == 0 {
		Help()
		return nil
	}
	switch args[0] {
	case "create":
		return Create(args[1:])
	case "get":
		return Get(args[1:])
	case "list":
		return List(args[1:])
	case "delete":
		return Delete(args[1:])
	case "-h", "--help":
		Help()
		return nil
	default:
		return fmt.Errorf("неизвестная подкоманда %q; доступны: create, get, list, delete", args[0])
	}
}

// newFlagSet отключает собственный вывод пакета flag: сообщения об ошибках
// пользователь должен видеть в одном стиле и на одном языке, поэтому
// формируем их сами и возвращаем как error.
func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

// parse разбирает флаги подкоманды и заодно ловит лишние позиционные аргументы.
func parse(fs *flag.FlagSet, args []string) error {
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			Help()
			return errHelpShown
		}
		return fmt.Errorf("mws profile %s: %w", fs.Name(), err)
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("mws profile %s: лишние аргументы: %s", fs.Name(), strings.Join(fs.Args(), " "))
	}
	return nil
}

// errHelpShown означает «справка уже напечатана, выходим без ошибки».
var errHelpShown = errors.New("справка показана")

// IsHelpShown отличает запрос справки от настоящей ошибки: в этом случае
// процесс должен завершиться с нулевым кодом и ничего не писать в stderr.
func IsHelpShown(err error) bool { return errors.Is(err, errHelpShown) }

type requiredFlag struct{ name, value string }

// requireFlags перечисляет сразу все недостающие флаги, чтобы пользователь не
// подбирал их по одному, запуская команду несколько раз.
func requireFlags(flags ...requiredFlag) error {
	var missing []string
	for _, f := range flags {
		if f.value == "" {
			missing = append(missing, "--"+f.name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("не указаны обязательные флаги: %s", strings.Join(missing, ", "))
	}
	return nil
}
