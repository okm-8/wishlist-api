package logger

import "api/internal/model/log"

func Log(ctx Context, level log.Level, message string, labels ...*log.Label) {
	record := log.NewRecord(level, message, labels)

	if ctx.FilterLogRecord(record) {
		ctx.PrintLogRecord(record)
	}
}
