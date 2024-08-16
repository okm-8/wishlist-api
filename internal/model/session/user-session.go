package session

import (
	"api/internal/model/user"
	"time"
)

type UserSession struct {
	session *Session
	user    *user.User
}

func NewUserSession(_user *user.User, expiresAt time.Time) *UserSession {
	return &UserSession{
		session: New(expiresAt),
		user:    _user,
	}
}

func RestoreUserSession(session *Session, _user *user.User) *UserSession {
	return &UserSession{
		session: session,
		user:    _user,
	}
}

func (session *UserSession) Session() *Session {
	return session.session
}

func (session *UserSession) Id() Id {
	return session.session.Id()
}

func (session *UserSession) User() *user.User {
	return session.user
}

func (session *UserSession) ExpireAt() time.Time {
	return session.session.ExpireAt()
}

func (session *UserSession) Expired() bool {
	return session.session.Expired()
}
