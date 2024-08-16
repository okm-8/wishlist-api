package session

import "time"

type SignupSession struct {
	session *Session
	email   string
}

func NewSignupSession(email string, expiresAt time.Time) *SignupSession {
	return &SignupSession{
		session: New(expiresAt),
		email:   email,
	}
}

func RestoreSignupSession(session *Session, email string) *SignupSession {
	return &SignupSession{
		session: session,
		email:   email,
	}
}

func (session *SignupSession) Session() *Session {
	return session.session
}

func (session *SignupSession) Id() Id {
	return session.session.Id()
}

func (session *SignupSession) Email() string {
	return session.email
}

func (session *SignupSession) ExpiresAt() time.Time {
	return session.session.ExpireAt()
}

func (session *SignupSession) Expired() bool {
	return session.session.Expired()
}
