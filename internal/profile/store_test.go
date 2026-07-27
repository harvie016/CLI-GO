package profile

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAndGet(t *testing.T) {
	dir := t.TempDir()
	s := NewStoreIn(dir)

	if err := s.Create("test", Profile{User: "example", Project: "new-project"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Формат файла — часть контракта, проверяем побайтово.
	data, err := os.ReadFile(filepath.Join(dir, "test.yaml"))
	if err != nil {
		t.Fatalf("файл профиля не создан: %v", err)
	}
	const want = "user: example\nproject: new-project\n"
	if string(data) != want {
		t.Errorf("содержимое файла = %q, ожидалось %q", data, want)
	}

	p, err := s.Get("test")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if p.User != "example" || p.Project != "new-project" {
		t.Errorf("Get вернул %+v", p)
	}
}

func TestCreateDuplicate(t *testing.T) {
	s := NewStoreIn(t.TempDir())
	if err := s.Create("test", Profile{User: "a", Project: "b"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	err := s.Create("test", Profile{User: "other", Project: "other"})
	if !errors.Is(err, ErrExists) {
		t.Fatalf("ожидалась ErrExists, получено %v", err)
	}
	// Повторный create не должен затирать данные.
	p, err := s.Get("test")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if p.User != "a" {
		t.Errorf("профиль перезаписан: %+v", p)
	}
}

func TestGetAndDeleteMissing(t *testing.T) {
	s := NewStoreIn(t.TempDir())

	if _, err := s.Get("нет-такого"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get: ожидалась ErrNotFound, получено %v", err)
	}
	if err := s.Delete("нет-такого"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Delete: ожидалась ErrNotFound, получено %v", err)
	}
}

func TestDelete(t *testing.T) {
	dir := t.TempDir()
	s := NewStoreIn(dir)
	if err := s.Create("test", Profile{User: "a", Project: "b"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := s.Delete("test"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "test.yaml")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("файл профиля остался на диске: %v", err)
	}
}

func TestListSortedAndSkipsForeignFiles(t *testing.T) {
	dir := t.TempDir()
	s := NewStoreIn(dir)
	for _, name := range []string{"beta", "alpha"} {
		if err := s.Create(name, Profile{User: "u-" + name, Project: "p-" + name}); err != nil {
			t.Fatalf("Create(%s): %v", name, err)
		}
	}
	// Мусор рядом с профилями: битый YAML, посторонний YAML и не-YAML файл.
	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("broken.yaml", "\tэто не yaml: [")
	write("foreign.yaml", "some: value\n")
	write("notes.txt", "user: x\n")

	got, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("List вернул %d профилей: %+v", len(got), got)
	}
	if got[0].Name != "alpha" || got[1].Name != "beta" {
		t.Errorf("список не отсортирован: %+v", got)
	}
	if got[0].User != "u-alpha" || got[0].Project != "p-alpha" {
		t.Errorf("данные профиля прочитаны неверно: %+v", got[0])
	}
}

func TestListEmptyDir(t *testing.T) {
	got, err := NewStoreIn(t.TempDir()).List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("ожидался пустой список, получено %+v", got)
	}
}

func TestInvalidNamesRejected(t *testing.T) {
	dir := t.TempDir()
	s := NewStoreIn(dir)

	names := []string{"", "../escape", "sub/test", `sub\test`, ".hidden", "a..b"}
	for _, name := range names {
		if err := s.Create(name, Profile{User: "u", Project: "p"}); err == nil {
			t.Errorf("Create(%q) прошёл, хотя имя недопустимо", name)
		}
		if _, err := s.Get(name); err == nil {
			t.Errorf("Get(%q) прошёл, хотя имя недопустимо", name)
		}
		if err := s.Delete(name); err == nil {
			t.Errorf("Delete(%q) прошёл, хотя имя недопустимо", name)
		}
	}

	// Ничего из недопустимых имён не должно было оказаться на диске.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("в каталоге появились файлы: %v", entries)
	}
}
