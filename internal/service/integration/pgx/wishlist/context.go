package wishlist

import (
	"api/internal/service/integration/pgx/driver"
	"context"
)

type Context interface {
	DriverContext() driver.Context
	RuntimeContext() context.Context
}
