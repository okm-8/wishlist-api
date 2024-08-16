package user

import (
	"api/internal/model/log"
	systemContext "api/internal/system/context"
	userSystem "api/internal/system/user"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listSignupTokenCmd = &cobra.Command{
	Use:   "list-signup-token <email>",
	Short: "list signup tokens",
	Long:  "list signup tokens",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context().(*systemContext.Context).WithLabels(log.NewLabel("command", "user-signup-token"))
		email := args[0]

		_tokens, err := userSystem.ListSignupToken(ctx, email)

		if err != nil {
			return err
		}

		items := make([]pterm.BulletListItem, len(_tokens)*4)

		for index, _token := range _tokens {
			items[index*4] = pterm.BulletListItem{
				Level:  0,
				Text:   fmt.Sprintf("Session ID: %s", _token.Session().Id()),
				Bullet: " ",
			}
			items[index*4+1] = pterm.BulletListItem{
				Level:  1,
				Text:   fmt.Sprintf("Expires at: %s", _token.Session().ExpireAt().Format("2006-01-02 15:04:05")),
				Bullet: " ",
			}
			items[index*4+2] = pterm.BulletListItem{
				Level:  1,
				Text:   fmt.Sprintf("For email: %s", email),
				Bullet: " ",
			}
			items[index*4+3] = pterm.BulletListItem{
				Level:  1,
				Text:   fmt.Sprintf("Token: %s", _token.String()),
				Bullet: " ",
			}
		}

		err = pterm.DefaultBulletList.WithItems(items).Render()

		if err != nil {
			return err
		}

		return nil
	},
}

func ListSignupTokenCmd() *cobra.Command {
	return listSignupTokenCmd
}
