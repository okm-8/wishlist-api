package migrations

import (
	"api/internal/model/migration"
	"api/internal/service/integration/pgx/driver"
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
	"time"
)

var (
	//go:embed query/create.sql
	createSQL string
	//go:embed query/insert.sql
	insertSQL string
	//go:embed query/select-by-filename.sql
	selectByFilenameSQL string
	//go:embed query/select-all.sql
	selectAllSQL string
)

func create(ctx context.Context, executor driver.CommandExecutor) error {
	_, err := executor.Exec(ctx, createSQL)

	return err
}

func insert(ctx context.Context, executor driver.CommandExecutor, filename string) error {
	_, err := executor.Exec(ctx, insertSQL, pgx.NamedArgs{
		"filename": filename,
	})

	return err
}

func selectByFilename(
	ctx context.Context,
	executor driver.QueryExecutor,
	filename string,
) (*migration.Migration, error) {
	rows, err := executor.Query(ctx, selectByFilenameSQL, pgx.NamedArgs{
		"filename": filename,
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var id uint64
	var executedAtTimestamp int64

	err = rows.Scan(&id, &executedAtTimestamp)

	if err != nil {
		return nil, err
	}

	return migration.Restore(id, filename, time.Unix(executedAtTimestamp, 0)), nil
}

func selectAll(ctx context.Context, executor driver.QueryExecutor) ([]*migration.Migration, error) {
	rows, err := executor.Query(ctx, selectAllSQL)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var migrations []*migration.Migration

	for rows.Next() {
		var id uint64
		var filename string
		var executedAtTimestamp time.Time

		err = rows.Scan(&id, &filename, &executedAtTimestamp)

		if err != nil {
			return nil, err
		}

		migrations = append(
			migrations,
			migration.Restore(id, filename, executedAtTimestamp),
		)
	}

	return migrations, nil
}

func Init(ctx Context) error {
	return driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		if err := create(ctx.RuntimeContext(), conn); err != nil {
			return err
		}

		return nil
	})
}

func FindByFilename(ctx Context, filename string) (*migration.Migration, error) {
	var _migration *migration.Migration

	err := driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		var err error
		_migration, err = selectByFilename(ctx.RuntimeContext(), conn, filename)

		return err
	})

	return _migration, err
}

func FindAll(ctx Context) ([]*migration.Migration, error) {
	var migrations []*migration.Migration

	err := driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		var err error
		migrations, err = selectAll(ctx.RuntimeContext(), conn)

		return err
	})

	return migrations, err
}

func ExecuteMigrations(ctx Context, migrations []*migration.Migration) error {
	path := ctx.MigrationsDirPath()

	return driver.Session(ctx.DriverContext(), func(conn *pgx.Conn) error {
		for _, _migration := range migrations {
			sql, err := os.ReadFile(fmt.Sprintf("%s/%s", path, _migration.Filename()))

			if err != nil {
				return err
			}

			fmt.Printf("Executing migration %s\n", _migration.Filename())
			fmt.Println(string(sql))

			_, err = conn.Exec(ctx.RuntimeContext(), string(sql))

			if err != nil {
				return err
			}

			err = insert(ctx.RuntimeContext(), conn, _migration.Filename())

			if err != nil {
				return err
			}
		}

		return nil
	})
}
