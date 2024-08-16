package user

import (
	"api/internal/model/session"
	"api/internal/service/integration/pgx/driver"
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5"
	"time"
)

var (
	//go:embed query/insert.sql
	insertSQL string
	//go:embed query/select-all.sql
	selectAllSQL string
	//go:embed query/delete-by-ids.sql
	deleteByIdsSQL string
)

func insert(ctx context.Context, executor driver.CommandExecutor, _session *session.Session) error {
	_, err := executor.Exec(ctx, insertSQL, pgx.NamedArgs{
		"id":        _session.Id().String(),
		"createdAt": _session.CreatedAt(),
		"expireAt":  _session.ExpireAt(),
	})

	return err
}

func selectAll(ctx context.Context, executor driver.QueryExecutor) ([]*session.Session, error) {
	rows, err := executor.Query(ctx, selectAllSQL)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions := make([]*session.Session, 0)

	for rows.Next() {
		var sessionIdStr string
		var sessionCreatedAt time.Time
		var sessionExpireAt time.Time

		err = rows.Scan(
			&sessionIdStr,
			&sessionCreatedAt,
			&sessionExpireAt,
		)

		if err != nil {
			return nil, err
		}

		sessionId, err := session.ParseId(sessionIdStr)

		if err != nil {
			return nil, err
		}

		_session := session.Restore(sessionId, sessionCreatedAt, sessionExpireAt)

		sessions = append(sessions, _session)
	}

	return sessions, nil
}

func del(ctx context.Context, executor driver.CommandExecutor, ids []session.Id) error {
	idsStr := make([]string, len(ids))

	for index, id := range ids {
		idsStr[index] = id.String()
	}

	_, err := executor.Exec(ctx, deleteByIdsSQL, pgx.NamedArgs{
		"ids": idsStr,
	})

	return err
}

func FindAll(ctx Context) ([]*session.Session, error) {
	var sessions []*session.Session
	var err error

	err = driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		sessions, err = selectAll(ctx.RuntimeContext(), conn)
		return err
	})

	return sessions, err
}

func Delete(ctx Context, ids ...session.Id) error {
	return driver.Transaction(ctx.DriverContext(), func(tx pgx.Tx) error {
		if err := del(ctx.RuntimeContext(), tx, ids); err != nil {
			return err
		}

		return nil
	}, driver.TxDeleteDefaultOpts)
}
