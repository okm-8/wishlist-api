package http

import (
	"net/http"
)

type MethodMap map[string]http.HandlerFunc

func Method(handlers MethodMap) http.HandlerFunc {
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
