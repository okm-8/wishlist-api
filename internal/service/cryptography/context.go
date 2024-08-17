package cryptography

import "api/internal/model/log"

type Context interface {
	Secret() string
	Log(level log.Level, message string, labels ...*log.Label)
}
