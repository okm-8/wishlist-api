package http

import (
	"api/internal/model/pagination"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func ReadPagination(request *http.Request, defaultValue *pagination.Pagination) (*pagination.Pagination, error) {
	pageStr := request.URL.Query().Get("page")
	limitStr := request.URL.Query().Get("limit")

	if pageStr == "" && limitStr == "" {
		return defaultValue, nil
	}

	var page, limit uint64
	var err error

	if pageStr == "" {
		page = defaultValue.Page()
	} else {
		page, err = strconv.ParseUint(pageStr, 10, 64)

		if err != nil {
			return nil, err
		}
	}

	if limitStr == "" {
		limit = defaultValue.Limit()
	} else {
		limit, err = strconv.ParseUint(limitStr, 10, 64)

		if err != nil {
			return nil, err
		}
	}

	return pagination.New(page, limit), nil
}

func ReadJson(request *http.Request, requestDto interface{}) error {
	err := json.NewDecoder(request.Body).Decode(requestDto)

	if err != nil {
		return err
	}

	return nil
}

const authHeader = "Authorization"

var ErrAuthHeaderInvalid = errors.New("authorization header is invalid")

func ReadAuthHeader(request *http.Request) (string, error) {
	authHeaderValue := request.Header.Get(authHeader)

	if authHeaderValue == "" || len(authHeaderValue) < 7 || authHeaderValue[:7] != "Bearer " {
		return "", ErrAuthHeaderInvalid
	}

	return authHeaderValue[7:], nil
}
