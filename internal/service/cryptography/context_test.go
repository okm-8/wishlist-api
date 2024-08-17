package cryptography

import "api/internal/model/log"

type contextStub struct {
	secret string
}

func newContextStub(secret string) *contextStub {
	return &contextStub{secret: secret}
}

func (c *contextStub) Secret() string {
	return c.secret
}

func (c *contextStub) Log(_ log.Level, _ string, _ ...*log.Label) {
	// Do nothing
}
