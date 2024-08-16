package driver

import "context"

type Context interface {
	RedisDsn() string
	RuntimeContext() context.Context
}
