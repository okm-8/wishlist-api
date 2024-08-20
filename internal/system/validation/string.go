package validation

import (
	"errors"
	"fmt"
	"regexp"
)

type StringRule func(value string) []error

func String(rules ...StringRule) Validator {
	return func(value any) []error {
		strValue, ok := value.(string)

		if !ok {
			return []error{errors.New("value is not a string")}
		}

		errs := make([]error, 0)

		for _, rule := range rules {
			if err := rule(strValue); err != nil {
				errs = append(errs, err...)
			}
		}

		return errs
	}
}

func NotEmpty() StringRule {
	return func(value string) []error {
		if value == "" {
			return []error{errors.New("value is empty")}
		}

		return nil
	}
}

func MinLength(length int) StringRule {
	return func(value string) []error {
		if len(value) < length {
			return []error{fmt.Errorf("value should be at least %d characters long", length)}
		}

		return nil
	}
}

func MaxLength(length int) StringRule {
	return func(value string) []error {
		if len(value) > length {
			return []error{fmt.Errorf("value should be at most %d characters long", length)}
		}

		return nil
	}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func Email() StringRule {
	return func(value string) []error {
		if !emailRegex.MatchString(value) {
			return []error{errors.New("value is not a valid email")}
		}

		return nil
	}
}
