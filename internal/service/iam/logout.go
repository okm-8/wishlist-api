package iam

import "api/internal/model/session"

func Logout(ctx Context, session *session.UserSession) error {
	if session == nil {
		return nil
	}

	return ctx.RevokeUserSession(session)
}
