package token

import (
	"bytes"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	t.Run("should parse token", func(t *testing.T) {
		ctx := newContextStub()
		token := New(ctx, []byte("data"), time.Now().Add(time.Hour))

		tokenString := token.String()

		parsedToken, err := Parse(ctx, tokenString)

		if err != nil {
			t.Errorf("expected token to be parsed, got %v", err)
		}

		if !bytes.Equal(token.Payload(), parsedToken.Payload()) {
			t.Errorf("expected payload to be %s, got %s", token.Payload(), parsedToken.Payload())
		}

		if token.ExpireAt().Truncate(time.Second) != parsedToken.ExpireAt() {
			t.Errorf("expected expire at to be %v, got %v", token.ExpireAt(), parsedToken.ExpireAt())
		}
	})
}
