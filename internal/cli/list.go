package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"mws/internal/profile"
)

// List — mws profile list
func List(args []string) error {
	fs := newFlagSet("list")
	if err := parse(fs, args); err != nil {
		return err
	}

	store, err := profile.NewStore()
	if err != nil {
		return err
	}
	profiles, err := store.List()
	if err != nil {
		return err
	}
	if len(profiles) == 0 {
		fmt.Println("В текущей директории нет профилей.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ИМЯ\tПОЛЬЗОВАТЕЛЬ\tПРОЕКТ")
	for _, p := range profiles {
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.User, p.Project)
	}
	return w.Flush()
}
