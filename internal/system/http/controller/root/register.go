package root

import (
	internalHttp "api/internal/service/integration/http"
	httpContext "api/internal/system/http/context"
	systemHttpServer "api/internal/system/http/server"
	"net/http"
)

func Register(server *systemHttpServer.Server) {
	server.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context().(*httpContext.RequestContext)

		internalHttp.WriteErrorResponse(ctx, writer, http.StatusNotFound, "not found", nil, nil)
	})
}
