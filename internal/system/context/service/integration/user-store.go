package integration

import (
	"api/internal/service/integration/pgx/driver"
	"context"
)

type UserStoreContext struct {
	ctx Context
}

func NewUserStoreContext(ctx Context) *UserStoreContext {
	return &UserStoreContext{
		ctx: ctx,
	}
}

func (ctx *UserStoreContext) PostgresDsn() string {
	return ctx.ctx.Config().Pgx().PostgresDsn()
}

func (ctx *UserStoreContext) RuntimeContext() context.Context {
	return ctx.ctx.RuntimeContext()
}

func (ctx *UserStoreContext) DriverContext() driver.Context {
	return ctx
}
