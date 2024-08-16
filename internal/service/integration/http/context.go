package http

import "api/internal/model/log"

type Context interface {
	Log(level log.Level, message string, labels ...*log.Label)
}
