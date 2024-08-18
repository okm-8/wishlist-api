package migrations

import (
	"api/internal/model/log"
	"api/internal/model/migration"
	systemContext "api/internal/system/context"
	"api/internal/system/migrations"
	_ "embed"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute [filename...]",
	Short: "Execute a migration",
	Long:  "Execute a migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		autoConfirm, err := cmd.Flags().GetBool("confirm")

		if err != nil {
			return err
		}

		force, err := cmd.Flags().GetBool("force")

		if err != nil {
			return err
		}

		ctx := cmd.Context().(*systemContext.Context).WithLabels(log.NewLabel("command", "migrations-execute"))

		all := len(args) == 0

		_migrations, err := migrations.List(ctx, all)

		if err != nil {
			return err
		}

		if len(args) != 0 {
			_migrationsMap := make(map[string]*migration.Migration)

			for _, _migration := range _migrations {
				_migrationsMap[_migration.Filename()] = _migration
			}

			filteredMigrations := make([]*migration.Migration, 0, len(args))

			for _, filename := range args {
				_migration, ok := _migrationsMap[filename]

				if !ok {
					pterm.FgRed.Printf("Migration %s not found\n", filename)

					return fmt.Errorf("migration %s not found", filename)
				}

				filteredMigrations = append(filteredMigrations, _migration)
			}

			_migrations = filteredMigrations
		}

		if len(_migrations) == 0 {
			pterm.FgGreen.Println("No migrations to execute")

			return nil
		}

		toExecute := make([]*migration.Migration, 0, len(_migrations))
		pterm.FgDefault.Println("Migrations to be executed:")
		list := make([]pterm.BulletListItem, 0, len(_migrations))

		for index, _migration := range _migrations {
			if _migration.ExecutedAt() == nil || force {
				toExecute = append(toExecute, _migration)
				list = append(list, pterm.BulletListItem{
					Level:       0,
					Text:        _migration.Filename(),
					TextStyle:   pterm.NewStyle(pterm.FgGreen),
					Bullet:      fmt.Sprintf("%d", index+1),
					BulletStyle: pterm.NewStyle(pterm.FgGreen),
				})
			} else {
				list = append(list, pterm.BulletListItem{
					Level:       0,
					Text:        fmt.Sprintf("%s (already executed)", _migration.Filename()),
					TextStyle:   pterm.NewStyle(pterm.FgGray),
					Bullet:      fmt.Sprintf("%d", index+1),
					BulletStyle: pterm.NewStyle(pterm.FgGray),
				})
			}
		}

		if !autoConfirm {
			err = pterm.DefaultBulletList.WithItems(list).Render()

			if err != nil {
				return err
			}

			confirm := pterm.DefaultInteractiveConfirm

			execute, err := confirm.Show("Do you want to execute these migrations?")

			if err != nil {
				return err
			}

			if !execute {
				return nil
			}
		}

		return migrations.Execute(ctx, toExecute)
	},
}

func ExecuteCmd() *cobra.Command {
	executeCmd.PersistentFlags().BoolP("confirm", "y", false, "Confirm the execution of migrations")
	executeCmd.Flags().BoolP("force", "f", false, "force re-execute executed migrations")

	return executeCmd
}
