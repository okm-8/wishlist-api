package logger

import "api/internal/model/log"

type Context interface {
	FilterLogRecord(record log.Record) bool
	PrintLogRecord(record log.Record)
}
