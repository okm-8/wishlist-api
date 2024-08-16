package environment

import "reflect"

func setString(value string, field *reflect.Value) error {
	field.SetString(value)

	return nil
}
