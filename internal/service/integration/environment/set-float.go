package environment

import (
	"reflect"
	"strconv"
)

func setFloat(value string, field *reflect.Value) error {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	field.SetFloat(floatValue)

	return nil
}
