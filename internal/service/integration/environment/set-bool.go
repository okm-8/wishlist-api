package environment

import (
	"reflect"
	"strconv"
)

func setBool(value string, field *reflect.Value) error {
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}

	field.SetBool(boolValue)

	return nil
}
