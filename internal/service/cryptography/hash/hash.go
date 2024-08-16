package hash

import (
	"crypto/hmac"
	"golang.org/x/crypto/blake2b"
)

func Hash(ctx Context, data []byte) []byte {
	_hash, err := blake2b.New256([]byte(ctx.Secret()))

	if err != nil {
		panic(err) // should never happen
	}

	_hash.Write(data)
	return _hash.Sum(nil)
}

func Verify(ctx Context, data []byte, signature []byte) bool {
	return hmac.Equal(signature, Hash(ctx, data))
}
