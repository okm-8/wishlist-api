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
	//go:embed query/signup-session/insert.sql
	insertSignupSessionSQL string
	//go:embed query/signup-session/select-by-id.sql
	selectSignupSessionByIdSQL string
	//go:embed query/signup-session/select-by-email.sql
	selectSignupSessionByEmailSQL string
)

func insertSignupSession(
	ctx context.Context,
	executor driver.CommandExecutor,
	session *session.SignupSession,
) error {
	err := insert(ctx, executor, session.Session())

	if err != nil {
		return err
	}

	_, err = executor.Exec(ctx, insertSignupSessionSQL, pgx.NamedArgs{
		"id":    session.Id().Bytes(),
		"email": session.Email(),
	})

	return err
}

func selectSignupSessionByEmail(
	ctx context.Context,
	executor driver.QueryExecutor,
	email string,
) ([]*session.SignupSession, error) {
	rows, err := executor.Query(ctx, selectSignupSessionByEmailSQL, pgx.NamedArgs{
		"email": email,
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions := make([]*session.SignupSession, 0)

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

		sessions = append(sessions, session.RestoreSignupSession(
			session.Restore(sessionId, sessionCreatedAt, sessionExpireAt),
			email,
		))
	}

	return sessions, nil
}

func selectSignupSessionById(
	ctx context.Context,
	executor driver.QueryExecutor,
	id session.Id,
) (*session.SignupSession, error) {
	rows, err := executor.Query(ctx, selectSignupSessionByIdSQL, pgx.NamedArgs{
		"id": id.String(),
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var email string
	var sessionCreatedAt time.Time
	var sessionExpireAt time.Time

	err = rows.Scan(
		&email,
		&sessionCreatedAt,
		&sessionExpireAt,
	)

	if err != nil {
		return nil, err
	}

	return session.RestoreSignupSession(
		session.Restore(id, sessionCreatedAt, sessionExpireAt),
		email,
	), nil
}

func StoreSignupSession(
	ctx Context,
	session *session.SignupSession,
) error {
	return driver.Transaction(ctx.DriverContext(), func(tx pgx.Tx) error {
		return insertSignupSession(ctx.RuntimeContext(), tx, session)
	}, driver.TxInsertDefaultOpts)
}

func FindSignupSessionByEmail(
	ctx Context,
	email string,
) ([]*session.SignupSession, error) {
	var sessions []*session.SignupSession
	var err error

	err = driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		sessions, err = selectSignupSessionByEmail(ctx.RuntimeContext(), conn, email)
		return err
	})

	return sessions, err
}

func FindSignupSessionById(
	ctx Context,
	id session.Id,
) (*session.SignupSession, error) {
	var _session *session.SignupSession
	var err error

	err = driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		_session, err = selectSignupSessionById(ctx.RuntimeContext(), conn, id)
		return err
	})

	return _session, err
}
