package user

import (
	"api/internal/model/log"
	userStore "api/internal/service/integration/pgx/user"
	systemContext "api/internal/system/context"
	userSystem "api/internal/system/user"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"time"
)

var signupTokenCmd = &cobra.Command{
	Use:   "signup-token <email>",
	Short: "Generate a signup token",
	Long:  "Generate a signup token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context().(*systemContext.Context).WithLabels(log.NewLabel("command", "user-signup-token"))
		email := args[0]
		expireInStr, err := cmd.Flags().GetString("expire-in")

		if err != nil {
			return err
		}

		expiresIn, err := time.ParseDuration(expireInStr)

		if err != nil {
			return err
		}

		_token, err := userSystem.SignupToken(ctx, email, expiresIn)

		if err != nil {
			if errors.Is(err, userStore.ErrEmailExists) {
				pterm.Error.Println("Email already exists")

				return nil
			}

			return err
		}

		items := []pterm.BulletListItem{
			{
				Level:  0,
				Text:   fmt.Sprintf("Session ID: %s", _token.Session().Id()),
				Bullet: " ",
			},
			{
				Level:  0,
				Text:   fmt.Sprintf("Expires at: %s", _token.Session().ExpireAt().Format("2006-01-02 15:04:05")),
				Bullet: " ",
			},
			{
				Level:  0,
				Text:   fmt.Sprintf("For email: %s", email),
				Bullet: " ",
			},
			{
				Level:  0,
				Text:   fmt.Sprintf("Token: %s", _token.String()),
				Bullet: " ",
			},
		}

		err = pterm.DefaultBulletList.WithItems(items).Render()

		if err != nil {
			return err
		}

		return nil
	},
}

func SignupTokenCmd() *cobra.Command {
	signupTokenCmd.Flags().StringP("expire-in", "e", "24h", "Token expiration time string")

	return signupTokenCmd
}
