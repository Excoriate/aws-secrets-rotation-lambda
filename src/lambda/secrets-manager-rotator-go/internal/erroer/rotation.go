package erroer

import "fmt"

type RotationError struct {
	Details string
	Err     error
}

func (e *RotationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Rotation has failed: %s: %s", e.Details, e.Err.Error())
	}
	return fmt.Sprintf("Rotator has failed: %s", e.Details)
}

func NewRotationError(details string, err error) *RotationError {
	return &RotationError{
		Details: fmt.Sprintf("Rotation has failed - %s", details),
		Err:     err,
	}
}
