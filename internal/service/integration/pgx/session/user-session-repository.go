package user

import (
	"api/internal/model/session"
	"api/internal/model/user"
	"api/internal/service/integration/pgx/driver"
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5"
	"time"
)

var (
	//go:embed query/user-session/insert.sql
	insertUserSessionSQL string
	//go:embed query/user-session/select-by-id.sql
	selectUserSessionByIdSQL string
)

func insertUserSession(
	ctx context.Context,
	executor driver.CommandExecutor,
	_userSession *session.UserSession,
) error {
	err := insert(ctx, executor, _userSession.Session())

	if err != nil {
		return err
	}

	_, err = executor.Exec(ctx, insertUserSessionSQL, pgx.NamedArgs{
		"id":     _userSession.Id().Bytes(),
		"userId": _userSession.User().Id().String(),
	})

	return err
}

func selectUserSessionById(
	ctx context.Context,
	executor driver.QueryExecutor,
	id session.Id,
) (*session.UserSession, error) {
	rows, err := executor.Query(ctx, selectUserSessionByIdSQL, pgx.NamedArgs{
		"id": id.String(),
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var userIdStr string
	var userEmail string
	var userName string
	var userAdmin bool
	var userPasswordHash []byte

	var sessionCreatedAt time.Time
	var sessionExpireAt time.Time

	err = rows.Scan(
		&userIdStr, &userEmail, &userName, &userAdmin, &userPasswordHash,
		&sessionCreatedAt, &sessionExpireAt,
	)

	if err != nil {
		return nil, err
	}

	userId, err := user.ParseId(userIdStr)

	if err != nil {
		return nil, err
	}

	return session.RestoreUserSession(
		session.Restore(
			id,
			sessionCreatedAt,
			sessionExpireAt,
		),
		user.Restore(userId, userEmail, userName, userAdmin, userPasswordHash),
	), nil
}

func StoreUserSession(ctx Context, userSession *session.UserSession) error {
	return driver.Transaction(ctx.DriverContext(), func(tx pgx.Tx) error {
		return insertUserSession(ctx.RuntimeContext(), tx, userSession)
	}, driver.TxInsertDefaultOpts)
}

func FindUserSession(ctx Context, id session.Id) (*session.UserSession, error) {
	var userSession *session.UserSession

	err := driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		var err error
		userSession, err = selectUserSessionById(ctx.RuntimeContext(), conn, id)

		return err
	})

	return userSession, err
}
