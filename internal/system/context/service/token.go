package service

import (
	"api/internal/model/log"
	"api/internal/service/cryptography"
	"errors"
)

var (
	ErrTokenHashInvalid = errors.New("token is invalid")
)

type TokenContext struct {
	ctx             Context
	cryptographyCtx cryptography.Context
}

func NewTokenContext(ctx Context) *TokenContext {
	return &TokenContext{
		ctx:             ctx,
		cryptographyCtx: NewCryptographyContext(ctx),
	}
}

func (ctx *TokenContext) Hash(payload []byte) []byte {
	return cryptography.Hash(ctx.cryptographyCtx, payload)
}

func (ctx *TokenContext) VerifyHash(payload []byte, _hash []byte) error {
	if cryptography.VerifyHash(ctx.cryptographyCtx, payload, _hash) {
		return nil
	}

	return ErrTokenHashInvalid
}

func (ctx *TokenContext) Log(level log.Level, message string, labels ...*log.Label) {
	ctx.ctx.Log(level, message, labels...)
}
