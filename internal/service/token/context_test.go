package token

import (
	"api/internal/model/log"
	"bytes"
)

type contextStub struct {
}

func newContextStub() *contextStub {
	return &contextStub{}
}

func (c *contextStub) Hash(_bytes []byte) []byte {
	return _bytes
}

func (c *contextStub) VerifyHash(_bytes []byte, hash []byte) error {
	if !bytes.Equal(_bytes, hash) {
		return ErrInvalid
	}

	return nil
}

func (c *contextStub) Log(_ log.Level, _ string, _ ...*log.Label) {
	// Do nothing
}
