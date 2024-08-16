package iam

import (
	"api/internal/model/user"
	"errors"
	"testing"
)

type credentials struct {
	email    string
	password string
}

func (creds *credentials) Email() string {
	return creds.email
}

func (creds *credentials) Password() string {
	return creds.password
}

func TestLogin(t *testing.T) {
	t.Run("should return session", func(t *testing.T) {
		ctx := newTestContext()
		sampleUser := user.New("test@test.com", "test")
		sampleUser.ChangePasswordHash([]byte("password"))
		_ = ctx.StoreUser(sampleUser)

		_session, err := Login(
			ctx,
			&credentials{
				email:    "test@test.com",
				password: "password",
			},
		)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if _session == nil {
			t.Errorf("expected session, got nil")
		}

		if _session.User().Id() != sampleUser.Id() {
			t.Errorf("unexpected user id: %s", _session.User().Id())
		}
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		ctx := newTestContext()

		_, err := Login(
			ctx,
			&credentials{
				email:    "test@test.com",
				password: "password",
			},
		)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !errors.Is(err, errUserNotFound) {
			t.Errorf("Expected error %v, got %v", errUserNotFound, err)
		}
	})

	t.Run("should return error if password is invalid", func(t *testing.T) {
		ctx := newTestContext()
		sampleUser := user.New("test@test.com", "test")
		sampleUser.ChangePasswordHash([]byte("password"))
		_ = ctx.StoreUser(sampleUser)

		_, err := Login(
			ctx,
			&credentials{
				email:    "test@test.com",
				password: "invalid",
			},
		)

		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
