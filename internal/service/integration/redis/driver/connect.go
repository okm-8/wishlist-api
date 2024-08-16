package driver

import (
	"github.com/redis/go-redis/v9"
)

func connect(ctx Context) (*redis.Client, error) {
	opts, err := redis.ParseURL(ctx.RedisDsn())

	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}
