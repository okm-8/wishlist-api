package hash

type contextStub struct {
	secret string
}

func newContextStub(secret string) *contextStub {
	return &contextStub{secret: secret}
}

func (c *contextStub) Secret() string {
	return c.secret
}
