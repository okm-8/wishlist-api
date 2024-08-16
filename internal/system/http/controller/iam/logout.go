package iam

import (
	"api/internal/service/iam"
	internalHttp "api/internal/service/integration/http"
	"api/internal/service/token"
	"errors"
	"net/http"
)

func isHeaderInvalid(err error) bool {
	return errors.Is(err, internalHttp.ErrAuthHeaderInvalid)
}

func isTokenInvalid(err error) bool {
	return errors.Is(err, token.ErrInvalid)
}

func isExpectedError(err error) bool {
	return isHeaderInvalid(err) || isTokenInvalid(err)
}

func logout(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)

	_session, err := ctx.RequestContext().UserSession()

	if err != nil && !isExpectedError(err) {
		internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)

		return
	}

	if err = iam.Logout(ctx.IamContext(), _session); err != nil {
		internalHttp.WriteInternalServerErrorResponse(ctx, writer, err)

		return
	}

	internalHttp.WriteSuccessResponse(ctx, writer, http.StatusOK, nil)
}
