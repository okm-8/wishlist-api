package iam

import (
	"api/internal/model/session"
	"api/internal/model/user"
	"testing"
	"time"
)

func TestLogout(t *testing.T) {
	t.Run("should logout", func(t *testing.T) {
		ctx := newTestContext()
		sampleUser := user.New("test@test.com", "test")
		sampleSession := session.NewUserSession(sampleUser, time.Now().Add(time.Hour))

		err := Logout(ctx, sampleSession)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}
