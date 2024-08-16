package token

import (
	"api/internal/model/log"
)

type Context interface {
	Hash(payload []byte) []byte
	VerifyHash(payload []byte, hash []byte) error
	Log(level log.Level, message string, labels ...*log.Label)
}
