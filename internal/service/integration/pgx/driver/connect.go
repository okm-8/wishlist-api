package driver

import (
	"github.com/jackc/pgx/v5"
)

func connect(ctx Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx.RuntimeContext(), ctx.PostgresDsn())
}
