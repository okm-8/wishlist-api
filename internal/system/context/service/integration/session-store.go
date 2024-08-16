package integration

import (
	"api/internal/service/integration/pgx/driver"
	"context"
)

type SessionStoreContext struct {
	ctx Context
}

func NewSessionStoreContext(ctx Context) *SessionStoreContext {
	return &SessionStoreContext{
		ctx: ctx,
	}
}

func (ctx *SessionStoreContext) PostgresDsn() string {
	return ctx.ctx.Config().Pgx().PostgresDsn()
}

func (ctx *SessionStoreContext) RuntimeContext() context.Context {
	return ctx.ctx.RuntimeContext()
}

func (ctx *SessionStoreContext) DriverContext() driver.Context {
	return ctx
}
