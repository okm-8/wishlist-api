package validation

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

func Optional(rules ...Validator) Validator {
	return func(value any) []error {
		ptrValue, ok := value.(*any)

		checkedValue := value
		if ok {
			if ptrValue == nil {
				return nil
			} else {
				checkedValue = *ptrValue
			}
		}

		return All(rules...)(checkedValue)
	}
}
