package integration

import (
	"api/internal/service/integration/pgx/driver"
	"context"
)

type WishlistStoreContext struct {
	ctx Context
}

func NewWishlistStoreContext(ctx Context) *WishlistStoreContext {
	return &WishlistStoreContext{
		ctx: ctx,
	}
}

func (ctx *WishlistStoreContext) PostgresDsn() string {
	return ctx.ctx.Config().Pgx().PostgresDsn()
}

func (ctx *WishlistStoreContext) RuntimeContext() context.Context {
	return ctx.ctx.RuntimeContext()
}

func (ctx *WishlistStoreContext) DriverContext() driver.Context {
	return ctx
}
