package http

import (
	"api/internal/model/log"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strings"
)

type responseStruct struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func writeHeaders(writer http.ResponseWriter, code int, headers map[string]string) {
	for key, value := range headers {
		writer.Header().Set(key, value)
	}
	writer.Header().Set("Content-Type", "application/json")

	writer.WriteHeader(code)
}

func writeBody(ctx Context, writer http.ResponseWriter, body []byte) {
	_, err := writer.Write(body)

	if err != nil {
		ctx.Log(log.Error, "failed to write response", log.NewLabel("error", err.Error()))
	}
}

func WriteSuccessResponse(ctx Context, writer http.ResponseWriter, code int, headers map[string]string) {
	if headers == nil {
		headers = make(map[string]string)
	}

	body, _ := json.Marshal(responseStruct{Success: true})

	writeHeaders(writer, code, headers)
	writeBody(ctx, writer, body)
}

func WriteErrorResponse(
	ctx Context,
	writer http.ResponseWriter,
	code int,
	message string,
	errs []error,
	headers map[string]string,
) {
	if headers == nil {
		headers = make(map[string]string)
	}

	errsStr := make([]string, len(errs))

	for i, err := range errs {
		errsStr[i] = err.Error()
	}

	body, _ := json.Marshal(responseStruct{Success: false, Message: message, Errors: errsStr})

	writeHeaders(writer, code, headers)
	writeBody(ctx, writer, body)
}

func WriteDataResponse(
	ctx Context,
	writer http.ResponseWriter,
	code int,
	payload any,
	headers map[string]string,
) {
	if headers == nil {
		headers = make(map[string]string)
	}

	body, err := json.Marshal(responseStruct{Success: true, Data: payload})

	if err != nil {
		WriteInternalServerErrorResponse(ctx, writer, err)
	}

	writeHeaders(writer, code, headers)
	writeBody(ctx, writer, body)
}

func WriteInfoResponse(
	ctx Context,
	writer http.ResponseWriter,
	code int,
	message string,
	headers map[string]string,
) {
	if headers == nil {
		headers = make(map[string]string)
	}

	body, _ := json.Marshal(responseStruct{Success: true, Message: message})

	writeHeaders(writer, code, headers)
	writeBody(ctx, writer, body)
}

func WriteInternalServerErrorResponse(ctx Context, writer http.ResponseWriter, err error) {
	ctx.Log(log.Error, err.Error(), log.NewLabel("debug", strings.Split(string(debug.Stack()), "\n")))

	WriteErrorResponse(
		ctx,
		writer,
		http.StatusInternalServerError,
		"internal server error",
		nil,
		nil,
	)
}
