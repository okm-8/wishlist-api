package service

import (
	"api/internal/model/log"
	jsonLogger "api/internal/service/integration/json-logger"
	"io"
)

type LoggerContext struct {
	ctx Context
}

func NewLoggerContext(ctx Context) *LoggerContext {
	return &LoggerContext{
		ctx: ctx,
	}
}

func (ctx *LoggerContext) FilterLogRecord(record log.Record) bool {
	return ctx.ctx.Config().LoggerConfig().MinLevel() <= record.Level()
}

func (ctx *LoggerContext) PrintLogRecord(record log.Record) {
	err := jsonLogger.PrintLogRecord(ctx, record)

	if err != nil {
		// if we cannot write to the log, we are in a bad state
		panic(err)
	}
}

func (ctx *LoggerContext) LogOutput() io.Writer {
	return ctx.ctx.Config().LoggerConfig().Output()
}

func (ctx *LoggerContext) ColoredLog() bool {
	return ctx.ctx.Config().LoggerConfig().Colored()
}
