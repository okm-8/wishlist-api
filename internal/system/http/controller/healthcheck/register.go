package healthcheck

import (
	internalHttp "api/internal/service/integration/http"
	"api/internal/system/http/server"
	"net/http"
)

func Register(server *server.Server) {
	server.HandleFunc("/healthcheck", internalHttp.Method(internalHttp.MethodMap{
		http.MethodGet: healthcheck,
	}))
}
