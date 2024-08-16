package context

import (
	"api/internal/model/log"
	"api/internal/model/session"
	"api/internal/model/user"
	internalHttp "api/internal/service/integration/http"
	sessionStore "api/internal/service/integration/pgx/session"
	"api/internal/service/token"
	systemContext "api/internal/system/context"
	"api/internal/system/context/service"
	"api/internal/system/context/service/integration"
	"errors"
	"net/http"
	"time"
)

var (
	ErrNotAUser = errors.New("not a user")
)

type RequestContext struct {
	request         *http.Request
	systemCtx       *systemContext.Context
	sessionStoreCtx sessionStore.Context
	tokenCtx        token.Context
	token           *token.SessionToken
	tokenErr        error
	userSession     *session.UserSession
	userSessionErr  error
}

func NewRequestContext(ctx *systemContext.Context, request *http.Request) *RequestContext {
	return &RequestContext{
		request:         request,
		systemCtx:       ctx,
		sessionStoreCtx: integration.NewSessionStoreContext(ctx),
		tokenCtx:        service.NewTokenContext(ctx),
	}
}

func (ctx *RequestContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.request.Context().Deadline()
}

func (ctx *RequestContext) Done() <-chan struct{} {
	return ctx.request.Context().Done()
}

func (ctx *RequestContext) Err() error {
	return ctx.request.Context().Err()
}

func (ctx *RequestContext) Value(key interface{}) interface{} {
	return ctx.request.Context().Value(key)
}

func (ctx *RequestContext) SystemContext() *systemContext.Context {
	return ctx.systemCtx
}

func (ctx *RequestContext) Token() (*token.SessionToken, error) {
	if ctx.tokenErr != nil {
		return nil, ctx.tokenErr
	}

	if ctx.token == nil {
		authHeader, err := internalHttp.ReadAuthHeader(ctx.request)

		if err != nil {
			ctx.tokenErr = err

			return nil, err
		}

		ctx.token, err = token.ParseSessionToken(ctx.tokenCtx, authHeader)

		if err != nil {
			ctx.tokenErr = err

			return nil, err
		}
	}

	return ctx.token, nil
}

func (ctx *RequestContext) UserSession() (*session.UserSession, error) {
	if ctx.userSessionErr != nil {
		return nil, ctx.userSessionErr
	}

	if ctx.userSession == nil {
		_token, err := ctx.Token()

		if err != nil {
			ctx.userSessionErr = err

			return nil, err
		}

		_userSession, err := sessionStore.FindUserSession(ctx.sessionStoreCtx, _token.Session().Id())

		if err != nil {
			ctx.userSessionErr = err

			return nil, err
		}

		if _userSession == nil {
			ctx.userSessionErr = ErrNotAUser

			return nil, ErrNotAUser
		}

		ctx.userSession = _userSession
	}

	return ctx.userSession, nil
}

func (ctx *RequestContext) Session() (*session.Session, error) {
	_token, err := ctx.Token()

	if err != nil {
		return nil, err
	}

	return _token.Session(), nil
}

func (ctx *RequestContext) User() (*user.User, error) {
	userSession, err := ctx.UserSession()

	if err != nil {
		return nil, err
	}

	return userSession.User(), nil
}

func (ctx *RequestContext) Log(level log.Level, message string, labels ...*log.Label) {
	ctx.systemCtx.Log(level, message, labels...)
}

func (ctx *RequestContext) WithLabels(labels ...*log.Label) *RequestContext {
	newCtx := ctx.systemCtx.WithLabels(labels...)

	return &RequestContext{
		request:         ctx.request,
		systemCtx:       newCtx,
		tokenCtx:        service.NewTokenContext(newCtx),
		sessionStoreCtx: integration.NewSessionStoreContext(newCtx),
		token:           ctx.token,
		tokenErr:        ctx.tokenErr,
		userSession:     ctx.userSession,
	}
}
