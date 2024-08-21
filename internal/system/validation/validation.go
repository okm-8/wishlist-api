package validation

import "fmt"

type Validator func(value any) []error

func All(rules ...Validator) Validator {
	return func(value any) []error {
		errs := make([]error, 0)

		for _, rule := range rules {
			if err := rule(value); err != nil {
				errs = append(errs, err...)
			}
		}

		return errs
	}
}

func Optional[T any](rules ...Validator) Validator {
	var t T

	return func(value any) []error {
		typedValue, ok := value.(*T)

		if !ok {
			return []error{fmt.Errorf("value is not a pointer to %T", t)}
		}

		if typedValue == nil {
			return nil
		}

		return All(rules...)(*typedValue)
	}
}
