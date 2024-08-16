package driver

import (
	"errors"
	"github.com/redis/go-redis/v9"
)

func Session(ctx Context, do func(conn *redis.Client) error) (err error) {
	conn, err := connect(ctx)

	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, conn.Close())
	}()

	return do(conn)
}

func Subscription(ctx Context, do func(sub *redis.PubSub) error, channels ...string) (err error) {
	conn, err := connect(ctx)

	if err != nil {
		return err
	}

	sub := conn.Subscribe(ctx.RuntimeContext(), channels...)

	defer func() {
		err = errors.Join(err, sub.Close())
	}()

	return do(sub)
}
