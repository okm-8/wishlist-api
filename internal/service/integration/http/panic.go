package http

import (
	"fmt"
	"net/http"
)

func Panic(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)

				if !ok {
					err = fmt.Errorf("%v", r)
				}

				WriteInternalServerErrorResponse(request.Context().(Context), writer, err)
			}
		}()

		handler(writer, request)
	}
}
