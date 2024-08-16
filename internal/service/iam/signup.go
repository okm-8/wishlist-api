package iam

import (
	"api/internal/model/session"
	"api/internal/model/user"
)

func SignUp(ctx Context, userInfo UserInfo) (*session.UserSession, error) {
	newUser := user.New(userInfo.Email(), userInfo.Name())
	ctx.SetUserPassword(newUser, userInfo.Password())

	err := ctx.StoreUser(newUser)

	if err != nil {
		return nil, err
	}

	return ctx.IssueUserSession(newUser)
}
