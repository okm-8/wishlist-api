package integration

import (
	"api/internal/service/integration/pgx/driver"
	"context"
)

type MigrationStoreContext struct {
	ctx Context
}

func NewMigrationStoreContext(ctx Context) *MigrationStoreContext {
	return &MigrationStoreContext{
		ctx: ctx,
	}
}

func (ctx *MigrationStoreContext) RuntimeContext() context.Context {
	return ctx.ctx.RuntimeContext()
}

func (ctx *MigrationStoreContext) PostgresDsn() string {
	return ctx.ctx.Config().Pgx().PostgresDsn()
}

func (ctx *MigrationStoreContext) DriverContext() driver.Context {
	return ctx
}

func (ctx *MigrationStoreContext) MigrationsDirPath() string {
	return ctx.ctx.Config().Migrations().DirPath()
}
