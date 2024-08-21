package iam

import (
	"api/internal/model/log"
	"api/internal/model/session"
	"api/internal/model/user"
	"api/internal/service/iam"
	sessionStore "api/internal/service/integration/pgx/session"
	userStore "api/internal/service/integration/pgx/user"
	"api/internal/service/token"
	systemContext "api/internal/system/context"
	"api/internal/system/context/service"
	"api/internal/system/context/service/integration"
	httpContext "api/internal/system/http/context"
	"net/http"
)

type Context struct {
	request         *http.Request
	requestCtx      *httpContext.RequestContext
	systemCtx       *systemContext.Context
	sessionStoreCtx sessionStore.Context
	tokenCtx        token.Context
	iamCtx          iam.Context
	userStoreCtx    userStore.Context
}

func NewContext(request *http.Request) *Context {
	requestCtx := request.Context().(*httpContext.RequestContext).WithLabels(log.NewLabel("controller", "iam"))
	systemCtx := requestCtx.SystemContext()

	return &Context{
		request:         request,
		requestCtx:      requestCtx,
		systemCtx:       systemCtx,
		sessionStoreCtx: integration.NewSessionStoreContext(systemCtx),
		tokenCtx:        service.NewTokenContext(systemCtx),
		iamCtx:          service.NewIamContext(systemCtx),
		userStoreCtx:    integration.NewUserStoreContext(systemCtx),
	}
}

func (ctx *Context) RequestContext() *httpContext.RequestContext {
	return ctx.requestCtx
}

func (ctx *Context) IamContext() iam.Context {
	return ctx.iamCtx
}

func (ctx *Context) User() *user.User {
	_user, err := ctx.requestCtx.User()

	if err != nil {
		panic(err) // This should never happen
	}

	return _user
}

func (ctx *Context) UserStoreContext() userStore.Context {
	return ctx.userStoreCtx
}

func (ctx *Context) UserSessionToken(_session *session.UserSession) *token.SessionToken {
	return token.NewSessionToken(ctx.tokenCtx, _session.Session())
}

func (ctx *Context) DecodeSessionTokenString(signInToken string) (*token.SessionToken, error) {
	return token.ParseSessionToken(ctx.tokenCtx, signInToken)
}

func (ctx *Context) FindSignupSession(id session.Id) (*session.SignupSession, error) {
	return sessionStore.FindSignupSessionById(ctx.sessionStoreCtx, id)
}

func (ctx *Context) CloseSession(_session *session.Session) error {
	return sessionStore.Delete(ctx.sessionStoreCtx, _session.Id())
}

func (ctx *Context) Log(level log.Level, message string, labels ...*log.Label) {
	ctx.systemCtx.Log(level, message, labels...)
}
