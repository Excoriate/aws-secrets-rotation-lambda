package erroer

import "fmt"

type RotatorConfigurationError struct {
	Details string
	Err     error
}

func (e *RotatorConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Rotator lambda configuration error: %s: %s", e.Details, e.Err.Error())
	}
	return fmt.Sprintf("Rotator lambda configuration error: %s", e.Details)
}

func NewConfigurationError(details string, err error) *RotatorConfigurationError {
	return &RotatorConfigurationError{
		Details: fmt.Sprintf("Unable to rotate secret - %s", details),
		Err:     err,
	}
}
