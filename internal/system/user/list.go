package user

import (
	"api/internal/model/log"
	"api/internal/model/pagination"
	"api/internal/model/user"
	userStore "api/internal/service/integration/pgx/user"
	systemContext "api/internal/system/context"
	"api/internal/system/context/service/integration"
)

func List(ctx *systemContext.Context, pagination *pagination.Pagination) ([]*user.User, error) {
	ctx = ctx.WithLabels(log.NewLabel("system", "user"), log.NewLabel("operation", "list"))
	userStoreCtx := integration.NewUserStoreContext(ctx)
	users, err := userStore.FindAll(userStoreCtx, pagination)

	if err != nil {
		return nil, err
	}

	return users, nil
}
