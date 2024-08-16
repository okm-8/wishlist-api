package iam

import (
	"api/internal/model/session"
	"api/internal/model/user"
	"bytes"
	"errors"
	"time"
)

var (
	errUserNotFound = errors.New("user not found")
	errUserExists   = errors.New("user already exists")
)

type testContext struct {
	user *user.User
}

func newTestContext() *testContext {
	return &testContext{
		user: nil,
	}
}

func (ctx *testContext) GetUserByEmail(email string) (*user.User, error) {
	if ctx.user != nil && ctx.user.Email() == email {
		return ctx.user, nil
	}

	return nil, errUserNotFound
}

func (ctx *testContext) IssueUserSession(_user *user.User) (*session.UserSession, error) {
	if _user != nil {
		return session.NewUserSession(_user, time.Now().Add(time.Hour)), nil
	}

	return nil, errUserNotFound
}

func (ctx *testContext) StoreUser(_user *user.User) error {
	if ctx.user != nil {
		return errUserExists
	}

	ctx.user = _user

	return nil
}

func (ctx *testContext) RevokeUserSession(_ *session.UserSession) error {
	return nil
}

func (ctx *testContext) SetUserPassword(_user *user.User, password string) {
	_user.ChangePasswordHash([]byte(password))
}

func (ctx *testContext) CheckUserPassword(_user *user.User, password string) bool {
	return bytes.Equal(_user.PasswordHash(), []byte(password))
}
