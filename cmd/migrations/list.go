package migrations

import (
	"api/internal/model/log"
	systemContext "api/internal/system/context"
	"api/internal/system/migrations"
	_ "embed"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List migrations",
	Long:  "List migrations",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		printAll, err := cmd.Flags().GetBool("all")

		if err != nil {
			return err
		}

		ctx := cmd.Context().(*systemContext.Context).WithLabels(log.NewLabel("command", "migrations-list"))
		_migrations, err := migrations.List(ctx, !printAll)

		if err != nil {
			return err
		}

		tableData := pterm.TableData{
			{"#", "Filename", "Executed at"},
		}

		count := len(_migrations)
		for index, migration := range _migrations {
			executedAt := pterm.FgRed.Sprint("Pending")
			if migration.ExecutedAt() != nil {
				executedAt = migration.ExecutedAt().Format("2006-01-02 15:04:05")
			}

			tableData = append(tableData, []string{
				fmt.Sprintf("%d", count-index),
				migration.Filename(),
				executedAt,
			})
		}

		var table string

		if len(tableData) == 1 {
			table = pterm.FgRed.Sprint("No migrations found")
		} else {
			table, err = pterm.DefaultTable.WithHasHeader().WithData(tableData).Srender()

			if err != nil {
				return err
			}

			table = table[:len(table)-1] // Remove the last newline
		}

		title := "Pending migrations"

		if printAll {
			title = "All migrations"
		}

		pterm.DefaultBox.WithTitle(title).Println(table)

		return nil
	},
}

func ListCmd() *cobra.Command {
	listCmd.PersistentFlags().BoolP("all", "a", false, "List all migrations")

	return listCmd
}
