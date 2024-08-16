package service

import (
	"api/internal/service/cryptography/hash"
)

func Hash(ctx Context, data []byte) []byte {
	return hash.Hash(ctx.Config().Cryptography(), data)
}

func VerifyHash(ctx Context, data []byte, _signature []byte) bool {
	return hash.Verify(ctx.Config().Cryptography(), data, _signature)
}
