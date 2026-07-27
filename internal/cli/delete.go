package cli

import (
	"fmt"

	"mws/internal/profile"
)

// Delete — mws profile delete --name=<name>
func Delete(args []string) error {
	fs := newFlagSet("delete")
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
	if err := store.Delete(*name); err != nil {
		return err
	}

	fmt.Printf("Профиль %s удалён\n", *name)
	return nil
}
