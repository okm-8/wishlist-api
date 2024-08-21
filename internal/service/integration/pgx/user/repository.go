package user

import (
	"api/internal/model/pagination"
	"api/internal/model/user"
	"api/internal/service/integration/pgx/driver"
	"context"
	_ "embed"
	"errors"
	"github.com/jackc/pgx/v5"
)

var (
	//go:embed query/insert.sql
	insertSQL string
	//go:embed query/update.sql
	updateSQL string
	//go:embed query/select-by-id.sql
	selectByIdSQL string
	//go:embed query/select-by-ids.sql
	selectByIdsSQL string
	//go:embed query/select-by-email.sql
	selectByEmailSQL string
	//go:embed query/select-all.sql
	selectAllSQL string
	//go:embed query/delete.sql
	deleteSQL string
)

var (
	ErrIdExists    = errors.New("id exists")
	ErrEmailExists = errors.New("email exists")
)

func insert(ctx context.Context, executor driver.CommandExecutor, user *user.User) error {
	_, err := executor.Exec(ctx, insertSQL, pgx.NamedArgs{
		"id":            user.Id().String(),
		"email":         user.Email(),
		"name":          user.Name(),
		"admin":         user.IsAdmin(),
		"password_hash": user.PasswordHash(),
	})

	return err
}

func update(ctx context.Context, executor driver.CommandExecutor, user *user.User) error {
	_, err := executor.Exec(ctx, updateSQL, pgx.NamedArgs{
		"id":           user.Id().String(),
		"email":        user.Email(),
		"name":         user.Name(),
		"admin":        user.IsAdmin(),
		"passwordHash": user.PasswordHash(),
	})

	return err
}

func selectById(ctx context.Context, executor driver.QueryExecutor, id user.Id) (*user.User, error) {
	rows, err := executor.Query(ctx, selectByIdSQL, pgx.NamedArgs{
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
	var name string
	var admin bool
	var passwordHash []byte

	err = rows.Scan(&email, &name, &admin, &passwordHash)

	if err != nil {
		return nil, err
	}

	return user.Restore(id, email, name, admin, passwordHash), nil
}

func selectByIds(ctx context.Context, executor driver.QueryExecutor, ids []user.Id) ([]*user.User, error) {
	idsStr := make([]string, 0, len(ids))

	for _, id := range ids {
		idsStr = append(idsStr, id.String())
	}

	rows, err := executor.Query(ctx, selectByIdsSQL, pgx.NamedArgs{
		"ids": idsStr,
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*user.User

	for rows.Next() {
		var idBytes []byte
		var email string
		var name string
		var admin bool
		var passwordHash []byte

		err = rows.Scan(&idBytes, &email, &name, &admin, &passwordHash)

		if err != nil {
			return nil, err
		}

		users = append(users, user.Restore(user.RestoreId(idBytes), email, name, admin, passwordHash))
	}

	return users, nil
}

func selectByEmail(ctx context.Context, executor driver.QueryExecutor, email string) (*user.User, error) {
	rows, err := executor.Query(ctx, selectByEmailSQL, pgx.NamedArgs{
		"email": email,
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var idBytes []byte
	var name string
	var admin bool
	var passwordHash []byte

	err = rows.Scan(&idBytes, &name, &admin, &passwordHash)

	if err != nil {
		return nil, err
	}

	return user.Restore(user.RestoreId(idBytes), email, name, admin, passwordHash), nil
}

func selectAll(
	ctx context.Context,
	executor driver.QueryExecutor,
	_pagination *pagination.Pagination,
) ([]*user.User, error) {
	rows, err := executor.Query(ctx, selectAllSQL, pgx.NamedArgs{
		"limit":  _pagination.Limit(),
		"offset": _pagination.Offset(),
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*user.User

	for rows.Next() {
		var idBytes []byte
		var email string
		var name string
		var admin bool
		var passwordHash []byte

		err = rows.Scan(&idBytes, &email, &name, &admin, &passwordHash)

		if err != nil {
			return nil, err
		}

		users = append(users, user.Restore(user.RestoreId(idBytes), email, name, admin, passwordHash))
	}

	return users, nil
}

func del(ctx context.Context, executor driver.CommandExecutor, id user.Id) error {
	_, err := executor.Exec(ctx, deleteSQL, pgx.NamedArgs{
		"id": id.Bytes(),
	})

	return err
}

func Store(ctx Context, user *user.User) error {
	return driver.Transaction(ctx.DriverContext(), func(tx pgx.Tx) error {
		userById, err := selectById(ctx.RuntimeContext(), tx, user.Id())

		if err != nil {
			return err
		}

		if userById != nil {
			return ErrIdExists
		}

		userByEmail, err := selectByEmail(ctx.RuntimeContext(), tx, user.Email())

		if err != nil {
			return err
		}

		if userByEmail != nil {
			return ErrEmailExists
		}

		if err := insert(ctx.RuntimeContext(), tx, user); err != nil {
			return err
		}

		return nil
	}, driver.TxInsertDefaultOpts)
}

func Update(ctx Context, id user.Id, doUpdate func(user *user.User) (*user.User, error)) error {
	return driver.Transaction(ctx.DriverContext(), func(tx pgx.Tx) error {
		_user, err := selectById(ctx.RuntimeContext(), tx, id)

		if err != nil {
			return err
		}

		_user, err = doUpdate(_user)

		if err != nil {
			return err
		}

		if _user == nil {
			return nil
		}

		userByEmail, err := selectByEmail(ctx.RuntimeContext(), tx, _user.Email())

		if err != nil {
			return err
		}

		if userByEmail != nil && userByEmail.Id() != _user.Id() {
			return ErrEmailExists
		}

		if err := update(ctx.RuntimeContext(), tx, _user); err != nil {
			return err
		}

		return nil
	}, driver.TxUpdateDefaultOpts)
}

func FindById(ctx Context, id user.Id) (*user.User, error) {
	var _user *user.User

	err := driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		var err error
		_user, err = selectById(ctx.RuntimeContext(), conn, id)

		return err
	})

	return _user, err
}

func FindByEmail(ctx Context, email string) (*user.User, error) {
	var _user *user.User

	err := driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		var err error
		_user, err = selectByEmail(ctx.RuntimeContext(), conn, email)

		return err
	})

	return _user, err
}

func FindAll(ctx Context, _pagination *pagination.Pagination) ([]*user.User, error) {
	var users []*user.User

	err := driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		var err error
		users, err = selectAll(ctx.RuntimeContext(), conn, _pagination)

		return err
	})

	return users, err
}

func FindAllByIds(ctx Context, ids []user.Id) ([]*user.User, error) {
	var users []*user.User

	err := driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		var err error
		users, err = selectByIds(ctx.RuntimeContext(), conn, ids)

		return err
	})

	return users, err
}

func Delete(ctx Context, id user.Id) error {
	return driver.Transaction(ctx.DriverContext(), func(tx pgx.Tx) error {
		return del(ctx.RuntimeContext(), tx, id)
	}, driver.TxDeleteDefaultOpts)
}
