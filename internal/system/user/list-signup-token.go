package user

import (
	"api/internal/model/log"
	sessionStore "api/internal/service/integration/pgx/session"
	"api/internal/service/token"
	systemContext "api/internal/system/context"
	"api/internal/system/context/service"
	"api/internal/system/context/service/integration"
)

func ListSignupToken(ctx *systemContext.Context, email string) ([]*token.SessionToken, error) {
	ctx = ctx.WithLabels(log.NewLabel("system", "user"), log.NewLabel("operation", "list-signup-token"))
	sessionStoreCtx := integration.NewSessionStoreContext(ctx)
	tokenCtx := service.NewTokenContext(ctx)

	sessions, err := sessionStore.FindSignupSessionByEmail(sessionStoreCtx, email)

	if err != nil {
		return nil, err
	}

	tokens := make([]*token.SessionToken, len(sessions))

	for i, _session := range sessions {
		tokens[i] = token.NewSessionToken(tokenCtx, _session.Session())
	}

	return tokens, nil
}
