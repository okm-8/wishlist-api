package session

import (
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	t.Run("should be expired", func(t *testing.T) {
		t.Parallel()

		session := New(time.Now().Add(-time.Hour))

		if !session.Expired() {
			t.Errorf("expected session to be expired")
		}
	})
}
