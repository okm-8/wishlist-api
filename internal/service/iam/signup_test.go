package iam

import (
	"api/internal/model/user"
	"bytes"
	"errors"
	"testing"
)

type userInfo struct {
	email    string
	name     string
	password string
}

func (data *userInfo) Email() string {
	return data.email
}

func (data *userInfo) Name() string {
	return data.name
}

func (data *userInfo) Password() string {
	return data.password
}

func TestSignUp(t *testing.T) {
	t.Run("should sign in", func(t *testing.T) {
		ctx := newTestContext()

		_session, err := SignUp(
			ctx,
			&userInfo{
				email:    "test@test.com",
				name:     "test",
				password: "password",
			},
		)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if _session.User().Email() != "test@test.com" {
			t.Errorf("unexpected email: %s", _session.User().Email())
		}

		if _session.User().Name() != "test" {
			t.Errorf("unexpected name: %s", _session.User().Name())
		}

		if !bytes.Equal(_session.User().PasswordHash(), []byte("password")) {
			t.Errorf("unexpected password: %s", _session.User().PasswordHash())
		}
	})

	t.Run("should fail to sign in if user already exists", func(t *testing.T) {
		ctx := newTestContext()
		sampleUser := user.New("test@test.com", "test")
		sampleUser.ChangePasswordHash([]byte("password"))
		_ = ctx.StoreUser(sampleUser)

		_, err := SignUp(
			ctx,
			&userInfo{
				email:    sampleUser.Email(),
				name:     sampleUser.Name(),
				password: "password",
			},
		)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !errors.Is(err, errUserExists) {
			t.Errorf("unexpected error: %s", err)
		}
	})
}
