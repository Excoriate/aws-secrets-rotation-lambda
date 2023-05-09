package erroer

import "fmt"

type SecretError struct {
	Details string
	Err     error
}

func (e *SecretError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Rotator lambda Secret error: %s: %s", e.Details, e.Err.Error())
	}
	return fmt.Sprintf("Rotator lambda secret error: %s", e.Details)
}

func NewSecretError(details string, err error) *SecretError {
	return &SecretError{
		Details: fmt.Sprintf("Unable to rotate secret - %s", details),
		Err:     err,
	}
}
