package integration

import (
	"api/internal/model/config"
	"context"
)

type Context interface {
	Config() *config.Config
	RuntimeContext() context.Context
}
