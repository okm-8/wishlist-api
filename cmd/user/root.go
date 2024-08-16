package user

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Long:  "Manage users",
}

func Root() *cobra.Command {
	rootCmd.AddCommand(ListCmd())
	rootCmd.AddCommand(SignupTokenCmd())
	rootCmd.AddCommand(ListSignupTokenCmd())

	return rootCmd
}
