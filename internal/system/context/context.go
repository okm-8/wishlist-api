package context

import (
	"api/internal/model/config"
	"api/internal/model/log"
	"api/internal/service/integration/output"
	"api/internal/service/logger"
	"api/internal/system/context/service"
	"api/internal/system/context/service/integration"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Context struct {
	ctx       context.Context
	cancel    context.CancelFunc
	env       *Environment
	config    *config.Config
	loggerCtx logger.Context
	labels    []*log.Label
}

func NewContext(logLabels ...*log.Label) (*Context, error) {
	env, err := NewEnvironment()

	if err != nil {
		return nil, err
	}

	logOutput, err := output.New(env.LogOutput)

	if err != nil {
		return nil, fmt.Errorf("failed to create log output: %w", err)
	}

	_config := config.NewConfig(
		config.NewCryptographyConfig(env.AppSecret),
		config.NewPgxConfig(env.PostgresDsn),
		config.NewMigrationsConfig(env.MigrationsPath),
		config.NewHttpConfig(
			config.NewHttpServerConfig(env.PublicAddress),
			config.NewHttpServerConfig(env.PrivateAddress),
		),
		config.NewLoggerConfig(
			log.ParseLevel(env.LogLevel),
			logOutput,
			env.ColoredLog,
		),
	)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	systemCtx := &Context{
		ctx:    ctx,
		cancel: cancel,
		config: _config,
		labels: logLabels,
	}

	go func() {
		<-sigs

		systemCtx.Cancel()
	}()

	systemCtx.loggerCtx = service.NewLoggerContext(systemCtx)

	return systemCtx, nil
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.ctx.Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

func (ctx *Context) Err() error {
	return ctx.ctx.Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.ctx.Value(key)
}

func (ctx *Context) RuntimeContext() context.Context {
	return ctx
}

func (ctx *Context) Config() *config.Config {
	return ctx.config
}

func (ctx *Context) IntegrationContext() integration.Context {
	return ctx
}

func (ctx *Context) Cancel() {
	ctx.cancel()
}

func (ctx *Context) Log(level log.Level, message string, labels ...*log.Label) {
	labels = append(labels, ctx.labels...)

	logger.Log(ctx.loggerCtx, level, message, labels...)
}

func (ctx *Context) WithLabels(labels ...*log.Label) *Context {
	newCtx := new(Context)
	*newCtx = *ctx

	newCtx.labels = append(newCtx.labels, labels...)

	return newCtx
}
