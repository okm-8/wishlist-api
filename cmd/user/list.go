package user

import (
	"api/internal/model/log"
	"api/internal/model/pagination"
	systemContext "api/internal/system/context"
	userSystem "api/internal/system/user"
	"errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Long:  "List users",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, err := cmd.Flags().GetUint64("limit")
		if err != nil {
			return err
		}

		if limit == 0 {
			return errors.New("limit must be greater than 0")
		}

		page, err := cmd.Flags().GetUint64("page")
		if err != nil {
			return err
		}

		if page == 0 {
			return errors.New("page must be greater than 0")
		}

		ctx := cmd.Context().(*systemContext.Context).WithLabels(log.NewLabel("command", "user-list"))

		_users, err := userSystem.List(ctx, pagination.New(page, limit))

		if err != nil {
			return err
		}

		tableData := pterm.TableData{
			{"ID", "Email", "Name", "Password set"},
		}

		for _, _user := range _users {
			passwordSet := pterm.FgRed.Sprint("X")

			if _user.PasswordSet() {
				passwordSet = pterm.FgGreen.Sprint("V")
			}

			tableData = append(tableData, []string{
				_user.Id().String(),
				_user.Email(),
				_user.Name(),
				passwordSet,
			})
		}

		var table string

		if len(tableData) == 1 {
			table = pterm.FgRed.Sprint("No users found")
		} else {
			table, err = pterm.DefaultTable.WithHasHeader().WithData(tableData).Srender()

			if err != nil {
				return err
			}

			table = table[:len(table)-1] // Remove the last newline
		}

		pterm.DefaultBox.WithTitle("Users").Println(table)

		return nil
	},
}

func ListCmd() *cobra.Command {
	listCmd.Flags().Uint64P("limit", "l", 10, "Limit the number of users to list")
	listCmd.Flags().Uint64P("page", "p", 1, "Page number to list")

	return listCmd
}
