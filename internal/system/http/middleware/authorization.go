package middleware

import (
	"api/internal/model/log"
	internalHttp "api/internal/service/integration/http"
	httpContext "api/internal/system/http/context"
	"errors"
	"net/http"
)

func AuthorizedOnly(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context().(*httpContext.RequestContext)

		user, err := ctx.User()

		if err != nil {
			ctx.Log(log.Debug, "failed to get user", log.NewLabel("error", err.Error()))

			internalHttp.WriteErrorResponse(
				ctx,
				writer,
				http.StatusUnauthorized,
				"unauthorized",
				[]error{errors.New("user is not authorized")},
				nil,
			)

			return
		}

		handler(writer, request.WithContext(ctx.WithLabels(
			log.NewLabel("userId", user.Id().String()),
		)))
	}
}

func AdminOnly(handler http.HandlerFunc) http.HandlerFunc {
	return AuthorizedOnly(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context().(*httpContext.RequestContext)

		user, err := ctx.User()

		if err != nil || !user.IsAdmin() {
			internalHttp.WriteErrorResponse(
				ctx,
				writer,
				http.StatusForbidden,
				"forbidden",
				[]error{errors.New("user is not an admin")},
				nil,
			)

			return
		}

		handler(writer, request)
	})
}
