package assert

import (
	"encoding/json"
	"testing"
)

func JsonEqualsValue[T comparable](t *testing.T, expectedValue T, jsonValue json.RawMessage, name string) {
	t.Helper()

	var value T

	err := json.Unmarshal(jsonValue, &value)

	if err != nil {
		t.Errorf("expected %s valid json, got %s: %v", name, jsonValue, err)
	}

	if value != expectedValue {
		t.Errorf("expected %s to be %v, got %v", name, expectedValue, value)
	}
}

func JsonNotEmpty(t *testing.T, jsonValue json.RawMessage, name string) {
	t.Helper()

	if len(jsonValue) == 0 {
		t.Errorf("expected %s non-empty json, got empty", name)
	}
}
