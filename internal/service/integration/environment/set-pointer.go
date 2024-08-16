package environment

import (
	"encoding/json"
	"reflect"
)

func setPointer(value string, field *reflect.Value) error {
	err := json.Unmarshal([]byte(value), field.Interface())

	if err != nil {
		return err
	}

	return nil
}
