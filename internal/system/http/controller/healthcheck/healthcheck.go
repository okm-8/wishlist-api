package healthcheck

import (
	"api/internal/model/log"
	internalHttp "api/internal/service/integration/http"
	"net/http"
)

func healthcheck(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request)

	ctx.Log(log.Info, "healthcheck")

	internalHttp.WriteSuccessResponse(ctx, writer, http.StatusOK, nil)
}
