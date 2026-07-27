package cli

import (
	"fmt"

	"mws/internal/profile"
)

// Create — mws profile create --name=<name> --user=<user> --project=<project>
func Create(args []string) error {
	fs := newFlagSet("create")
	name := fs.String("name", "", "имя профиля")
	user := fs.String("user", "", "имя пользователя")
	project := fs.String("project", "", "название проекта")
	if err := parse(fs, args); err != nil {
		return err
	}
	if err := requireFlags(
		requiredFlag{"name", *name},
		requiredFlag{"user", *user},
		requiredFlag{"project", *project},
	); err != nil {
		return err
	}

	store, err := profile.NewStore()
	if err != nil {
		return err
	}
	if err := store.Create(*name, profile.Profile{User: *user, Project: *project}); err != nil {
		return err
	}

	fmt.Printf("Профиль %s создан (%s.yaml)\n", *name, *name)
	return nil
}
