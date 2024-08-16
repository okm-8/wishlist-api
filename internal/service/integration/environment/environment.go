package environment

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"reflect"
	"slices"
)

var setters = map[reflect.Kind]func(string, *reflect.Value) error{
	reflect.String:  setString,
	reflect.Int:     setInt,
	reflect.Int8:    setInt,
	reflect.Int16:   setInt,
	reflect.Int32:   setInt,
	reflect.Int64:   setInt,
	reflect.Uint:    setUint,
	reflect.Uint8:   setUint,
	reflect.Uint16:  setUint,
	reflect.Uint32:  setUint,
	reflect.Uint64:  setUint,
	reflect.Float32: setFloat,
	reflect.Float64: setFloat,
	reflect.Bool:    setBool,
	reflect.Pointer: setPointer,
}

func Read(a interface{}, envFiles ...string) (err error) {
	// reverse the order of the files because dotenv does not override the values
	slices.Reverse(envFiles)

	err = godotenv.Load(envFiles...)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// it is ok if the file does not exist

		return err
	}

	aPointer := reflect.TypeOf(a)

	if aPointer.Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer to a struct, got %v", aPointer.Kind())
	}

	aType := aPointer.Elem()

	if aType.Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct, got a pointer to %v", aType.Elem().Kind())
	}

	aValue := reflect.ValueOf(a).Elem()
	errs := make([]error, 0)

	for i := 0; i < aType.NumField(); i++ {
		field := aType.Field(i)
		value := aValue.Field(i)
		tag := field.Tag.Get("env")
		defaultValue := field.Tag.Get("default")
		isRequired := field.Tag.Get("required") == "true"

		if tag == "" && isRequired {
			return fmt.Errorf("required field %v is not set", field.Name)
		}

		if tag != "" {
			setter, ok := setters[field.Type.Kind()]
			if !ok {
				return fmt.Errorf("unsupported type: %v", field.Type.Kind())
			}

			envValue := os.Getenv(tag)
			if envValue == "" {
				envValue = defaultValue
			}

			if envValue != "" {
				if err = setter(envValue, &value); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	return errors.Join(errs...)
}
