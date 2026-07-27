package cli

import (
	"fmt"

	"mws/internal/profile"
)

// Get — mws profile get --name=<name>
func Get(args []string) error {
	fs := newFlagSet("get")
	name := fs.String("name", "", "имя профиля")
	if err := parse(fs, args); err != nil {
		return err
	}
	if err := requireFlags(requiredFlag{"name", *name}); err != nil {
		return err
	}

	store, err := profile.NewStore()
	if err != nil {
		return err
	}
	p, err := store.Get(*name)
	if err != nil {
		return err
	}

	fmt.Printf("Профиль:      %s\n", *name)
	fmt.Printf("Пользователь: %s\n", p.User)
	fmt.Printf("Проект:       %s\n", p.Project)
	return nil
}
