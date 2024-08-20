package iam

import (
	"api/internal/service/iam"
	internalHttp "api/internal/service/integration/http"
	"api/internal/system/validation"
	"errors"
	"net/http"
)

type loginRequest struct {
	EmailValue    string `json:"email"`
	PasswordValue string `json:"password"`
}

func (request *loginRequest) Email() string {
	return request.EmailValue
}

func (request *loginRequest) Password() string {
	return request.PasswordValue
}

var validateLoginRequest = validation.Struct(
	validation.StructField(
		"email",
		func(request *loginRequest) any {
			return request.EmailValue
		},
		validation.String(
			validation.NotEmpty(),
			validation.Email(),
		),
	),
	validation.StructField(
		"password",
		func(request *loginRequest) any {
			return request.PasswordValue
		},
		validation.String(
			validation.NotEmpty(),
		),
	),
)

type loginResponse struct {
	Token string `json:"token"`
}

func login(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)
	_loginRequest := new(loginRequest)

	err := internalHttp.ReadJson(request, _loginRequest)

	if err != nil {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", []error{err}, nil)

		return
	}

	if errs := validateLoginRequest(_loginRequest); len(errs) > 0 {
		internalHttp.WriteErrorResponse(ctx, writer, http.StatusBadRequest, "invalid request", errs, nil)

		return
	}

	_session, err := iam.Login(ctx.IamContext(), _loginRequest)

	if err != nil {
		if errors.Is(err, iam.ErrUserNotFound) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "user not found", nil, nil)
		} else if errors.Is(err, iam.ErrInvalidCredentials) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusUnauthorized, "invalid password", nil, nil)
		} else if errors.Is(err, iam.ErrPasswordNotSet) {
			internalHttp.WriteErrorResponse(ctx, writer, http.StatusUnauthorized, "password not set", nil, nil)
		} else {
			internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)
		}

		return
	}

	internalHttp.WriteDataResponse(ctx, writer, http.StatusOK, loginResponse{
		Token: ctx.UserSessionToken(_session).String(),
	}, nil)
}
