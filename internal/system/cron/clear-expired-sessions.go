package cron

import (
	"api/internal/model/log"
	"api/internal/model/session"
	sessionStore "api/internal/service/integration/pgx/session"
	systemContext "api/internal/system/context"
	"api/internal/system/context/service/integration"
	"fmt"
	"time"
)

func clearExpiredSessions(ctx *systemContext.Context) error {
	sessionStoreCtx := integration.NewSessionStoreContext(ctx.IntegrationContext())
	sessions, err := sessionStore.FindAll(sessionStoreCtx)

	if err != nil {
		return err
	}

	toRemove := make([]session.Id, 0)

	for _, _session := range sessions {
		if _session.Expired() {
			toRemove = append(toRemove, _session.Id())
		}
	}

	toRemoveCount := len(toRemove)
	if toRemoveCount > 0 {
		ctx.Log(
			log.Info,
			fmt.Sprintf("clearing %d expired sessions", toRemoveCount),
		)

		if err := sessionStore.Delete(sessionStoreCtx, toRemove...); err != nil {
			ctx.Log(
				log.Error,
				"failed to clear expired sessions",
				log.NewLabel("error", err.Error()),
			)

			return err
		}
	}

	return nil
}

func ClearExpiredSessionsCron(ctx *systemContext.Context, period time.Duration) error {
	ctx = ctx.WithLabels(log.NewLabel("cron", "clear-expired-sessions"), log.NewLabel("period", period.String()))

	if err := clearExpiredSessions(ctx); err != nil {
		return err
	}

	ctx.Log(
		log.Info,
		fmt.Sprintf("clearing expired sessions every %s", period),
	)

	if err := cron(ctx, clearExpiredSessions, period); err != nil {
		ctx.Log(
			log.Error,
			"clearing expired sessions failed",
			log.NewLabel("error", err.Error()),
		)
	}

	ctx.Log(
		log.Info,
		"clearing expired sessions stopped",
	)

	return nil
}
