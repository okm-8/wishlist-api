package user

import (
	"api/internal/model/log"
	"api/internal/model/session"
	sessionStore "api/internal/service/integration/pgx/session"
	userStore "api/internal/service/integration/pgx/user"
	"api/internal/service/token"
	systemContext "api/internal/system/context"
	"api/internal/system/context/service"
	"api/internal/system/context/service/integration"
	"time"
)

func SignupToken(ctx *systemContext.Context, email string, expireIn time.Duration) (*token.SessionToken, error) {
	ctx = ctx.WithLabels(log.NewLabel("system", "user"), log.NewLabel("operation", "signup-token"))
	userStoreCtx := integration.NewUserStoreContext(ctx)
	sessionStoreCtx := integration.NewSessionStoreContext(ctx)
	tokenCtx := service.NewTokenContext(ctx)

	user, err := userStore.FindByEmail(userStoreCtx, email)

	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, userStore.ErrEmailExists
	}

	_session := session.NewSignupSession(email, time.Now().Add(expireIn))

	if err = sessionStore.StoreSignupSession(sessionStoreCtx, _session); err != nil {
		return nil, err
	}

	_token := token.NewSessionToken(tokenCtx, _session.Session())

	return _token, nil
}
