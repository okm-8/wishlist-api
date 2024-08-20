package validation

import "fmt"

type StructRule[T any] func(structValue T) []error

func StructField[T any](name string, getField func(structValue T) any, rules ...Validator) StructRule[T] {
	return func(structValue T) []error {
		fieldValue := getField(structValue)

		errs := make([]error, 0)

		for _, rule := range rules {
			if err := rule(fieldValue); err != nil {
				for _, err := range err {
					errs = append(errs, fmt.Errorf("%s: %w", name, err))
				}
			}
		}

		return errs
	}
}

func Struct[T any](rules ...StructRule[T]) Validator {
	var t T

	return func(value any) []error {
		structValue, ok := value.(T)

		if !ok {
			return []error{fmt.Errorf("value is not a %T", t)}
		}

		errs := make([]error, 0)

		for _, rule := range rules {
			if err := rule(structValue); err != nil {
				errs = append(errs, err...)
			}
		}

		return errs
	}
}
