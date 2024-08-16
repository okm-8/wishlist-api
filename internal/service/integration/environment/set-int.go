package environment

import (
	"reflect"
	"strconv"
)

func setInt(value string, field *reflect.Value) error {
	intValue, err := strconv.Atoi(value)

	if err != nil {
		return err
	}

	field.SetInt(int64(intValue))

	return nil
}

func setUint(value string, field *reflect.Value) error {
	uintValue, err := strconv.ParseUint(value, 10, 64)

	if err != nil {
		return err
	}

	field.SetUint(uintValue)

	return nil
}
