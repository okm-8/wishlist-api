package migrations

import (
	"github.com/spf13/cobra"
)

var migrationsCmd = &cobra.Command{
	Use:   "migrations",
	Short: "Migration is a command line tool for managing database migrations",
	Long:  "Migration is a command line tool for managing database migrations",
}

func Root() *cobra.Command {
	migrationsCmd.AddCommand(CreateCmd())
	migrationsCmd.AddCommand(ExecuteCmd())
	migrationsCmd.AddCommand(ListCmd())

	return migrationsCmd
}
