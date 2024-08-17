package cryptography

import (
	"api/internal/model/log"
	"crypto/hmac"
	"golang.org/x/crypto/blake2b"
)

func Hash(ctx Context, data []byte) []byte {
	_hash, err := blake2b.New256([]byte(ctx.Secret()))

	if err != nil {
		// should never happen
		ctx.Log(log.Error, "failed to create hash", log.NewLabel("error", err.Error()))

		panic(err)
	}

	_hash.Write(data)

	return _hash.Sum(nil)
}

func VerifyHash(ctx Context, data []byte, signature []byte) bool {
	return hmac.Equal(signature, Hash(ctx, data))
}
