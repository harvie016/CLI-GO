package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Store читает и пишет профили в каталоге dir.
type Store struct {
	dir string
}

// NewStore привязывает хранилище к текущей рабочей директории: по условию
// профили лежат только в той папке, откуда запущена команда.
func NewStore() (*Store, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("не удалось определить текущую директорию: %w", err)
	}
	return &Store{dir: dir}, nil
}

// NewStoreIn создаёт хранилище над произвольным каталогом — нужно тестам.
func NewStoreIn(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name+ext)
}

// Create записывает новый профиль и отказывается трогать уже существующий.
func (s *Store) Create(name string, p Profile) error {
	if err := validateName(name); err != nil {
		return err
	}
	data, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("не удалось сформировать YAML: %w", err)
	}
	// O_EXCL вместо связки Stat+WriteFile: проверка существования и создание
	// файла происходят одной операцией, промежутка для гонки не остаётся.
	f, err := os.OpenFile(s.path(name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%w: %s", ErrExists, name)
		}
		return fmt.Errorf("не удалось создать файл профиля: %w", err)
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("не удалось записать профиль: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("не удалось закрыть файл профиля: %w", err)
	}
	return nil
}

// Get возвращает профиль по имени.
func (s *Store) Get(name string) (Profile, error) {
	if err := validateName(name); err != nil {
		return Profile{}, err
	}
	return s.read(name)
}

// Delete удаляет файл профиля.
func (s *Store) Delete(name string) error {
	if err := validateName(name); err != nil {
		return err
	}
	if err := os.Remove(s.path(name)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w: %s", ErrNotFound, name)
		}
		return fmt.Errorf("не удалось удалить профиль: %w", err)
	}
	return nil
}

// List возвращает все профили каталога, отсортированные по имени.
// Посторонние и битые YAML-файлы пропускаются: список не должен падать из-за
// чужого файла, случайно оказавшегося рядом.
func (s *Store) List() ([]Named, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать текущую директорию: %w", err)
	}

	var profiles []Named
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ext {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ext)
		if validateName(name) != nil {
			continue
		}
		p, err := s.read(name)
		if err != nil {
			continue
		}
		// Оба поля пустые — это какой-то посторонний YAML, а не наш профиль:
		// create без непустых user и project профиль не создаёт.
		if p.User == "" && p.Project == "" {
			continue
		}
		profiles = append(profiles, Named{Name: name, Profile: p})
	}

	sort.Slice(profiles, func(i, j int) bool { return profiles[i].Name < profiles[j].Name })
	return profiles, nil
}

func (s *Store) read(name string) (Profile, error) {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Profile{}, fmt.Errorf("%w: %s", ErrNotFound, name)
		}
		return Profile{}, fmt.Errorf("не удалось прочитать файл профиля: %w", err)
	}
	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Profile{}, fmt.Errorf("файл профиля %s повреждён: %w", name+ext, err)
	}
	return p, nil
}
