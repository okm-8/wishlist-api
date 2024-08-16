package service

import (
	"api/internal/model/log"
	"errors"
)

var (
	ErrTokenHashInvalid = errors.New("token is invalid")
)

type TokenContext struct {
	ctx Context
}

func NewTokenContext(ctx Context) *TokenContext {
	return &TokenContext{
		ctx: ctx,
	}
}

func (ctx *TokenContext) Hash(payload []byte) []byte {
	return Hash(ctx.ctx, payload)
}

func (ctx *TokenContext) VerifyHash(payload []byte, hash []byte) error {
	if VerifyHash(ctx.ctx, payload, hash) {
		return nil
	}

	return ErrTokenHashInvalid
}

func (ctx *TokenContext) Log(level log.Level, message string, labels ...*log.Label) {
	ctx.ctx.Log(level, message, labels...)
}
