package errors

import "fmt"

type TaskError struct {
	Details string
	Err     error
}

func (e *TaskError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Task error: %s: %s", e.Details, e.Err.Error())
	}
	return fmt.Sprintf("Task error: %s", e.Details)
}

func NewTaskError(details string, err error) *TaskError {
	return &TaskError{
		Details: fmt.Sprintf("Unable to complete task %s", details),
		Err:     err,
	}
}
