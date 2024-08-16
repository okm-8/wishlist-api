package server

import (
	"api/internal/model/log"
	internalHttp "api/internal/service/integration/http"
	systemContext "api/internal/system/context"
	httpContext "api/internal/system/http/context"
	"github.com/google/uuid"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
	ctx *systemContext.Context
}

func New(ctx *systemContext.Context) *Server {
	return &Server{
		mux: http.NewServeMux(),
		ctx: ctx,
	}
}

func (server *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	server.mux.HandleFunc(pattern, handler)
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := httpContext.NewRequestContext(
		server.ctx.WithLabels(
			log.NewLabel("requestId", uuid.New().String()),
			log.NewLabel("uri", request.URL.Path),
		),
		request,
	)
	request = request.WithContext(ctx)

	internalHttp.Panic(server.mux.ServeHTTP)(writer, request)
}
