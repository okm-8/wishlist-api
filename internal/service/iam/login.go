package iam

import (
	"api/internal/model/session"
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordNotSet     = errors.New("password not set")
	ErrUserNotFound       = errors.New("user not found")
)

func Login(ctx Context, credentials Credentials) (*session.UserSession, error) {
	user, err := ctx.GetUserByEmail(credentials.Email())

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if !user.PasswordSet() {
		return nil, ErrPasswordNotSet
	}

	if !ctx.CheckUserPassword(user, credentials.Password()) {
		return nil, ErrInvalidCredentials
	}

	return ctx.IssueUserSession(user)
}
