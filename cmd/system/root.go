package system

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "system",
	Short: "Manage system",
	Long:  "Manage system",
}

func Root() *cobra.Command {
	rootCmd.AddCommand(Core())
	return rootCmd
}
