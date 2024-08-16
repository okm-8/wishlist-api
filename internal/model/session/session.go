package session

import (
	"time"
)

type Session struct {
	id        Id
	createdAt time.Time
	expireAt  time.Time
}

func New(expireAt time.Time) *Session {
	return &Session{
		id:        newId(),
		createdAt: time.Now(),
		expireAt:  expireAt,
	}
}

func Restore(id Id, createdAt time.Time, expireAt time.Time) *Session {
	return &Session{
		id:        id,
		createdAt: createdAt,
		expireAt:  expireAt,
	}
}

func (session *Session) Id() Id {
	return session.id
}

func (session *Session) CreatedAt() time.Time {
	return session.createdAt
}

func (session *Session) ExpireAt() time.Time {
	return session.expireAt
}

func (session *Session) Expired() bool {
	return session.expireAt.Before(time.Now())
}
