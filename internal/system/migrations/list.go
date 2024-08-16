package migrations

import (
	"api/internal/model/log"
	"api/internal/model/migration"
	migrationsStore "api/internal/service/integration/pgx/migrations"
	systemContext "api/internal/system/context"
	"api/internal/system/context/service/integration"
	"os"
	"strings"
)

func List(ctx *systemContext.Context, newOnly bool) ([]*migration.Migration, error) {
	ctx = ctx.WithLabels(log.NewLabel("system", "migrations"), log.NewLabel("operation", "list"))
	storeCtx := integration.NewMigrationStoreContext(ctx)

	err := migrationsStore.Init(storeCtx)

	if err != nil {
		return nil, err
	}

	_migrations, err := migrationsStore.FindAll(storeCtx)

	migrationsIndex := make(map[string]*migration.Migration, len(_migrations))

	for _, _migration := range _migrations {
		migrationsIndex[_migration.Filename()] = _migration
	}

	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(ctx.Config().Migrations().DirPath())

	if err != nil {
		return nil, err
	}

	newMigrations := make([]*migration.Migration, 0)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), EXT) {
			continue
		}

		_, ok := migrationsIndex[entry.Name()]

		if !ok {
			newMigrations = append(newMigrations, migration.New(entry.Name()))
		}
	}

	if newOnly {
		return newMigrations, nil
	}

	return append(newMigrations, _migrations...), nil
}
