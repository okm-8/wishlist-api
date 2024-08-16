package driver

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	TxInsertDefaultOpts = pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	}
	TxUpdateDefaultOpts = pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}
	TxDeleteDefaultOpts = TxUpdateDefaultOpts
)

type QueryExecutor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type CommandExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func Session(ctx Context, do func(conn *pgx.Conn) error) (err error) {
	conn, err := connect(ctx)

	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, conn.Close(ctx.RuntimeContext()))
	}()

	return do(conn)
}

func Transaction(ctx Context, do func(tx pgx.Tx) error, opts pgx.TxOptions) error {
	return Session(ctx, func(conn *pgx.Conn) (err error) {
		tx, err := conn.BeginTx(ctx.RuntimeContext(), opts)

		if err != nil {
			return err
		}

		defer func() {
			if r := recover(); r != nil {
				_ = tx.Rollback(ctx.RuntimeContext())

				panic(r)
			}
		}()

		err = do(tx)

		if err != nil {
			return errors.Join(err, tx.Rollback(ctx.RuntimeContext()))
		}

		return tx.Commit(ctx.RuntimeContext())
	})
}
