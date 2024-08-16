package driver

import "context"

type Context interface {
	PostgresDsn() string
	RuntimeContext() context.Context
}
