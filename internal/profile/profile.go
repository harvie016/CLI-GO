// Package profile описывает профиль и правила его хранения в YAML-файлах.
package profile

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// Profile — содержимое файла профиля. Имя профиля внутри не хранится:
// оно определяется именем файла, профиль test лежит в test.yaml.
type Profile struct {
	User    string `yaml:"user"`
	Project string `yaml:"project"`
}

// Named — профиль вместе с именем, восстановленным из имени файла.
type Named struct {
	Name string
	Profile
}

// Ошибки, по которым CLI-слой отличает предсказуемые ситуации от реальных
// сбоев файловой системы. Сравнивать тексты сообщений не нужно — errors.Is.
var (
	ErrExists   = errors.New("профиль уже существует")
	ErrNotFound = errors.New("профиль не найден")
)

// ext — расширение файлов профилей.
const ext = ".yaml"

// validateName следит за тем, чтобы имя профиля осталось именем файла в текущей
// папке и не превратилось в путь куда-то ещё.
func validateName(name string) error {
	switch {
	case name == "":
		return errors.New("имя профиля не может быть пустым")
	case strings.ContainsAny(name, `/\`):
		return fmt.Errorf("имя профиля %q не должно содержать разделители пути", name)
	case strings.Contains(name, ".."):
		return fmt.Errorf("имя профиля %q не должно содержать \"..\"", name)
	case strings.HasPrefix(name, "."):
		return fmt.Errorf("имя профиля %q не должно начинаться с точки", name)
	case filepath.Base(name) != name:
		// Страховка от того, что не поймали проверки выше: например, от
		// windows-специфичных конструкций вроде "C:name".
		return fmt.Errorf("имя профиля %q не может быть путём", name)
	}
	return nil
}
