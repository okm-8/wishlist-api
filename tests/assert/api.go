package assert

import (
	"encoding/json"
	"testing"
)

func apiResponse(
	t *testing.T,
	responsePayload json.RawMessage,
	expectedSuccess bool,
	expectedMessage string,
	expectedData interface{},
	expectedErrors bool,
) {
	t.Helper()

	_apiResponse := make(map[string]json.RawMessage)

	err := json.Unmarshal(responsePayload, &_apiResponse)

	if err != nil {
		t.Errorf("expected response to be valid json, got %s: %v", responsePayload, err)
	}

	JsonEqualsValue(t, expectedSuccess, _apiResponse["success"], "success")

	if expectedMessage != "" {
		JsonEqualsValue(t, expectedMessage, _apiResponse["message"], "message")
	}

	if expectedErrors {
		JsonNotEmpty(t, _apiResponse["errors"], "errors")
	}

	if expectedData != nil {
		JsonNotEmpty(t, _apiResponse["data"], "data")

		err = json.Unmarshal(_apiResponse["data"], expectedData)

		if err != nil {
			t.Errorf("expected data to be valid json, got %s: %v", _apiResponse["data"], err)
		}
	}
}

func ApiSuccessResponse(t *testing.T, responsePayload json.RawMessage) {
	t.Helper()

	apiResponse(t, responsePayload, true, "", nil, false)
}

func ApiDataResponse(t *testing.T, responsePayload json.RawMessage, data interface{}) {
	t.Helper()

	apiResponse(t, responsePayload, true, "", data, false)
}

func ApiResponseError(t *testing.T, responsePayload json.RawMessage, message string, hasErrors bool) {
	t.Helper()

	apiResponse(t, responsePayload, false, message, nil, hasErrors)
}
