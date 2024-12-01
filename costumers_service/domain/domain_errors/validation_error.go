package domain_errors

import "fmt"

type ValidationError struct {
	Fields   map[string]error
	InnerErr error
}

func (e ValidationError) Error() string {
	if len(e.Fields) > 1 {
		return fmt.Sprintf("validation errors found: %v", e.InnerErr)
	}

	return fmt.Sprintf("validation error found : %v", e.InnerErr)
}

func (e *ValidationError) AddField(name string, err error) {
	e.Fields[name] = err
}

func NewValidationError(err error) ValidationError {
	return ValidationError{
		InnerErr: err,
		Fields:   make(map[string]error),
	}
}
