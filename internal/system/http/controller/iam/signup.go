package iam

import (
	"api/internal/model/log"
	"api/internal/service/iam"
	internalHttp "api/internal/service/integration/http"
	userStore "api/internal/service/integration/pgx/user"
	"api/internal/service/token"
	"errors"
	"net/http"
)

type signUpRequest struct {
	EmailValue    string `json:"email"`
	NameValue     string `json:"name"`
	PasswordValue string `json:"password"`
}

func (request *signUpRequest) Validate() []error {
	var result []error

	if request.EmailValue == "" {
		result = append(result, errors.New("email is required"))
	}

	if request.NameValue == "" {
		result = append(result, errors.New("name is required"))
	}

	if request.PasswordValue == "" {
		result = append(result, errors.New("password is required"))
	}

	return result
}

func (request *signUpRequest) Email() string {
	return request.EmailValue
}

func (request *signUpRequest) Name() string {
	return request.NameValue
}

func (request *signUpRequest) Password() string {
	return request.PasswordValue
}

type singUpResponse struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

const queryToken = "token"

func isSingUpTokenInvalid(err error) bool {
	return errors.Is(err, token.ErrInvalid)
}

func signUp(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)
	tokenStr := request.URL.Query().Get(queryToken)
	singUpToken, err := ctx.DecodeSessionTokenString(tokenStr)

	if err != nil {
		ctx.Log(log.Debug, "invalid token", log.NewLabel("error", err.Error()), log.NewLabel("token", tokenStr))

		if isSingUpTokenInvalid(err) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusUnauthorized, "invalid token", nil, nil)
		} else {
			internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)
		}

		return
	}

	if singUpToken.Session().Expired() {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusUnauthorized, "token expired", nil, nil)

		return
	}

	_signUpRequest := new(signUpRequest)

	err = internalHttp.ReadJson(request, _signUpRequest)

	if err != nil {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid JSON", nil, nil)

		return
	}

	if errs := _signUpRequest.Validate(); len(errs) > 0 {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", errs, nil)

		return
	}

	_signupSession, err := ctx.FindSignupSession(singUpToken.Session().Id())

	if err != nil {
		internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)

		return
	}

	if _signupSession == nil {
		ctx.Log(log.Warning, "session not found")
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusUnauthorized, "invalid token", nil, nil)

		return
	}

	if _signupSession.Expired() {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusUnauthorized, "token expired", nil, nil)

		return
	}

	if _signupSession.Email() != _signUpRequest.Email() {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusUnauthorized, "invalid token", nil, nil)

		return
	}

	_session, err := iam.SignUp(ctx.IamContext(), _signUpRequest)

	if err != nil {
		if errors.Is(err, userStore.ErrEmailExists) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusConflict, "email already exists", nil, nil)
		} else {
			internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)
		}

		return
	}

	err = ctx.CloseSession(_signupSession.Session())

	if err != nil {
		internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusCreated, singUpResponse{
		UserId: _session.User().Id().String(),
		Token:  ctx.UserSessionToken(_session).String(),
	}, nil)
}

func signupNoToken(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)

	_signUpRequest := new(signUpRequest)

	err := internalHttp.ReadJson(request, _signUpRequest)

	if err != nil {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", []error{err}, nil)

		return
	}

	if errs := _signUpRequest.Validate(); len(errs) > 0 {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", errs, nil)

		return
	}

	_session, err := iam.SignUp(ctx.IamContext(), _signUpRequest)

	if err != nil {
		if errors.Is(err, userStore.ErrEmailExists) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusConflict, "email already exists", nil, nil)
		} else {
			internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)
		}

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusCreated, singUpResponse{
		UserId: _session.User().Id().String(),
		Token:  ctx.UserSessionToken(_session).String(),
	}, nil)
}
