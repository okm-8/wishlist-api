package migrations

import (
	"api/internal/model/log"
	"api/internal/model/migration"
	migrationsStore "api/internal/service/integration/pgx/migrations"
	systemContext "api/internal/system/context"
	"api/internal/system/context/service/integration"
)

func Execute(ctx *systemContext.Context, migrations []*migration.Migration) error {
	ctx = ctx.WithLabels(log.NewLabel("system", "migrations"), log.NewLabel("operation", "execute"))
	migrationsStoreCtx := integration.NewMigrationStoreContext(ctx)

	err := migrationsStore.Init(migrationsStoreCtx)

	if err != nil {
		return err
	}

	return migrationsStore.ExecuteMigrations(migrationsStoreCtx, migrations)
}
