package cron

import (
	systemContext "api/internal/system/context"
	"time"
)

func cron(ctx *systemContext.Context, operation func(*systemContext.Context) error, period time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(period):
			if err := operation(ctx); err != nil {
				return err
			}
		}
	}
}
