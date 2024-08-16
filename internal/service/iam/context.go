package iam

import (
	"api/internal/model/session"
	"api/internal/model/user"
)

type Credentials interface {
	Email() string
	Password() string
}

type UserInfo interface {
	Email() string
	Name() string
	Password() string
}

type Context interface {
	StoreUser(_user *user.User) error
	GetUserByEmail(email string) (*user.User, error)
	IssueUserSession(_user *user.User) (*session.UserSession, error)
	RevokeUserSession(session *session.UserSession) error
	SetUserPassword(_user *user.User, password string)
	CheckUserPassword(_user *user.User, password string) bool
}
