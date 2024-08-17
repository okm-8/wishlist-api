package service

import (
	"api/internal/model/log"
)

type CryptographyContext struct {
	ctx Context
}

func NewCryptographyContext(ctx Context) *CryptographyContext {
	return &CryptographyContext{ctx: ctx}
}

func (ctx *CryptographyContext) Secret() string {
	return ctx.ctx.Config().Cryptography().Secret()
}

func (ctx *CryptographyContext) Log(level log.Level, message string, labels ...*log.Label) {
	ctx.ctx.Log(level, message, labels...)
}
