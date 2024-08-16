package service

import (
	"api/internal/model/session"
	"api/internal/model/user"
	sessionStore "api/internal/service/integration/pgx/session"
	userStore "api/internal/service/integration/pgx/user"
	"api/internal/system/context/service/integration"
	"time"
)

type IamContext struct {
	ctx             Context
	userStoreCtx    userStore.Context
	sessionStoreCtx sessionStore.Context
}

func NewIamContext(ctx Context) *IamContext {
	return &IamContext{
		ctx:             ctx,
		userStoreCtx:    integration.NewUserStoreContext(ctx.IntegrationContext()),
		sessionStoreCtx: integration.NewSessionStoreContext(ctx.IntegrationContext()),
	}
}

func (ctx *IamContext) UserStoreContext() userStore.Context {
	return ctx.userStoreCtx
}

func (ctx *IamContext) SessionStoreContext() sessionStore.Context {
	return ctx.sessionStoreCtx
}

func (ctx *IamContext) StoreUser(_user *user.User) error {
	return userStore.Store(ctx.UserStoreContext(), _user)
}

func (ctx *IamContext) GetUserByEmail(email string) (*user.User, error) {
	return userStore.FindByEmail(ctx.UserStoreContext(), email)
}

func (ctx *IamContext) IssueUserSession(_user *user.User) (*session.UserSession, error) {
	_session := session.NewUserSession(_user, time.Now().Add(time.Hour))

	if err := sessionStore.StoreUserSession(ctx.UserStoreContext(), _session); err != nil {
		return nil, err
	}

	return _session, nil
}

func (ctx *IamContext) RevokeUserSession(session *session.UserSession) error {
	return sessionStore.Delete(ctx.SessionStoreContext(), session.Id())
}

func (ctx *IamContext) SetUserPassword(_user *user.User, password string) {
	hash := Hash(ctx.ctx, []byte(password))
	_user.ChangePasswordHash(hash)
}

func (ctx *IamContext) CheckUserPassword(_user *user.User, password string) bool {
	return VerifyHash(ctx.ctx, []byte(password), _user.PasswordHash())
}
