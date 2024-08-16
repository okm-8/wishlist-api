package token

import (
	"api/internal/model/session"
	"testing"
	"time"
)

func TestSessionToken(t *testing.T) {
	t.Run("should parse session token", func(t *testing.T) {
		ctx := newContextStub()
		_session := session.New(time.Now().Add(time.Hour))

		token := NewSessionToken(ctx, _session)

		tokenString := token.String()

		parsedToken, err := ParseSessionToken(ctx, tokenString)

		if err != nil {
			t.Errorf("expected session token to be parsed, got %v", err)
		}

		if parsedToken.Session().Id() != _session.Id() {
			t.Errorf("expected session id to be %s, got %s", _session.Id(), parsedToken.Session().Id())
		}

		if parsedToken.Session().CreatedAt() != _session.CreatedAt().Truncate(time.Second) {
			t.Errorf(
				"expected session created at to be %v, got %v",
				_session.CreatedAt(),
				parsedToken.Session().CreatedAt(),
			)
		}

		if parsedToken.Session().ExpireAt() != _session.ExpireAt().Truncate(time.Second) {
			t.Errorf(
				"expected session expire at to be %v, got %v",
				_session.ExpireAt().Truncate(time.Second),
				parsedToken.Session().ExpireAt(),
			)
		}
	})
}
