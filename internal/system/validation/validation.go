package validation

type Validator func(value any) []error

func All(value any, rules ...Validator) []error {
	errs := make([]error, 0)

	for _, rule := range rules {
		if err := rule(value); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}
