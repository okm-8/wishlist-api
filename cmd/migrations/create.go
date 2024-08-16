package migrations

import (
	"api/internal/model/log"
	systemContext "api/internal/system/context"
	"api/internal/system/migrations"
	_ "embed"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <description>",
	Short: "Create a new migration",
	Long:  "Create a new migration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context().(*systemContext.Context).WithLabels(log.NewLabel("command", "migrations-create"))

		filename, err := migrations.Create(ctx, args[0])

		if err != nil {
			return err
		}

		pterm.BgDefault.Printfln("Migration %s created", pterm.FgGreen.Sprint(filename))

		return nil
	},
}

func CreateCmd() *cobra.Command {
	return createCmd
}
