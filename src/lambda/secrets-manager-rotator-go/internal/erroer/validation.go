package erroer

import "fmt"

type RotatorValidationError struct {
	Details string
	Err     error
}

func (e *RotatorValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Rotator lambda validation error: %s: %s", e.Details, e.Err.Error())
	}
	return fmt.Sprintf("Rotator lambda validation error: %s", e.Details)
}

func NewValidationError(details string, err error) *RotatorValidationError {
	return &RotatorValidationError{
		Details: fmt.Sprintf("Unable to rotate secret - %s", details),
		Err:     err,
	}
}
