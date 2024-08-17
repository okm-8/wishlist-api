package token

import (
	"api/internal/model/session"
	"encoding/binary"
	"errors"
	"time"
)

type SessionToken struct {
	token   *Token
	session *session.Session
}

const sessionPayloadSize = 32

var ErrInvalidSessionPayload = errors.New("invalid session payload")

func serializeSession(session *session.Session) []byte {
	payload := session.Id().Bytes()
	payload = binary.BigEndian.AppendUint64(payload, uint64(session.CreatedAt().Unix()))
	payload = binary.BigEndian.AppendUint64(payload, uint64(session.ExpireAt().Unix()))
	return payload
}

func deserializeSession(payload []byte) (*session.Session, error) {
	if len(payload) < sessionPayloadSize {
		return nil, ErrInvalidSessionPayload
	}

	id := session.RestoreId(payload[:16])
	createdAt := time.Unix(int64(binary.BigEndian.Uint64(payload[16:24])), 0)
	expireAt := time.Unix(int64(binary.BigEndian.Uint64(payload[24:32])), 0)

	return session.Restore(id, createdAt, expireAt), nil
}

func NewSessionToken(ctx Context, _session *session.Session) *SessionToken {
	payload := serializeSession(_session)

	token := New(ctx, payload, _session.ExpireAt())

	return &SessionToken{
		token:   token,
		session: _session,
	}
}

func ParseSessionToken(ctx Context, tokenStr string) (*SessionToken, error) {
	token, err := Parse(ctx, tokenStr)
	if err != nil {
		return nil, err
	}

	_session, err := deserializeSession(token.Payload())

	if err != nil {
		return nil, err
	}

	return &SessionToken{
		token:   token,
		session: _session,
	}, nil
}

func (token *SessionToken) Session() *session.Session {
	return token.session
}

func (token *SessionToken) Expired() bool {
	return token.token.Expired()
}

func (token *SessionToken) String() string {
	return token.token.String()
}
