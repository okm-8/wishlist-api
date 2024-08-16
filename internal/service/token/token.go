package token

import (
	"api/internal/model/log"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

var ErrInvalid = errors.New("invalid token")

type Token struct {
	payload  []byte
	expireAt time.Time
	token    string
}

func New(ctx Context, payload []byte, expireAt time.Time) *Token {
	var expireAtBytes = make([]byte, 8)
	binary.BigEndian.PutUint64(expireAtBytes, uint64(expireAt.Unix()))

	fullPayload := append(expireAtBytes, payload...)

	var fullPayloadSize = make([]byte, 4)
	binary.BigEndian.PutUint32(fullPayloadSize, uint32(len(fullPayload)))

	tokenBytes := append(fullPayloadSize, fullPayload...)
	hash := ctx.Hash(tokenBytes)

	token := base64.URLEncoding.EncodeToString(append(tokenBytes, hash...))

	return &Token{
		payload:  payload,
		expireAt: expireAt,
		token:    token,
	}
}

func Parse(ctx Context, token string) (*Token, error) {
	tokenBytes, err := base64.URLEncoding.DecodeString(token)

	if err != nil {
		ctx.Log(log.Debug, "cannot encode token bytes", log.NewLabel("error", err.Error()))

		return nil, ErrInvalid
	}

	tokenBytesLen := len(tokenBytes)
	if tokenBytesLen < 12 {
		ctx.Log(log.Debug, "token is too short", log.NewLabel("length", tokenBytesLen))

		return nil, ErrInvalid
	}

	payloadSize := binary.BigEndian.Uint32(tokenBytes[:4])
	hash := tokenBytes[payloadSize+4:]
	fullPayload := tokenBytes[:payloadSize+4]

	if err = ctx.VerifyHash(fullPayload, hash); err != nil {
		fmt.Printf("payload: %s\n", fullPayload)
		ctx.Log(log.Debug, "token hash is invalid", log.NewLabel("error", err.Error()))

		return nil, ErrInvalid
	}

	fullPayload = fullPayload[4:]
	expireAt := time.Unix(int64(binary.BigEndian.Uint64(fullPayload[:8])), 0)
	payload := fullPayload[8:]

	return &Token{
		payload:  payload,
		expireAt: expireAt,
		token:    token,
	}, nil
}

func (t *Token) Payload() []byte {
	return t.payload
}

func (t *Token) ExpireAt() time.Time {
	return t.expireAt
}

func (t *Token) Expired() bool {
	return t.expireAt.Before(time.Now())
}

func (t *Token) String() string {
	return t.token
}
