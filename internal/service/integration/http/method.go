package http

import (
	"net/http"
)

func Method(handlers map[string]http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if handler, ok := handlers[request.Method]; ok {
			handler(writer, request)
		} else {
			WriteErrorResponse(
				request.Context().(Context),
				writer,
				http.StatusMethodNotAllowed,
				"method not allowed",
				nil,
				nil,
			)
		}
	}
}
