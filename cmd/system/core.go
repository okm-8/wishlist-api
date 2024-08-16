package system

import (
	"api/internal/model/log"
	systemContext "api/internal/system/context"
	"api/internal/system/cron"
	httpSystem "api/internal/system/http"
	"context"
	"errors"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

var coreCmd = &cobra.Command{
	Use:   "core",
	Short: "Run core system",
	Long:  "Run core system",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context().(*systemContext.Context).WithLabels(log.NewLabel("command", "core"))
		var group sync.WaitGroup
		errChan := make(chan error, 1)

		group.Add(1)
		go func() {
			defer group.Done()

			if err := cron.ClearExpiredSessionsCron(ctx, time.Hour); err != nil {
				errChan <- err
			}
		}()

		group.Add(1)
		go func() {
			defer group.Done()

			if err := httpSystem.PublicHttp(ctx); err != nil {
				errChan <- err
			}
		}()

		group.Add(1)
		go func() {
			defer group.Done()

			if err := httpSystem.PrivateHttp(ctx); err != nil {
				errChan <- err
			}
		}()

		select {
		case err := <-errChan:
			ctx.Log(log.Error, "core system failed", log.NewLabel("error", err.Error()))
			ctx.Cancel()
		case <-ctx.Done():
			ctx.Log(log.Info, "core system context is done")

			if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
				ctx.Log(log.Error, "core system context is done with error", log.NewLabel("error", err.Error()))
			}
		}

		group.Wait()

		return nil
	},
}

func Core() *cobra.Command {
	return coreCmd
}
