package llama

import "fmt"

type Error struct {
	Op      string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

func newError(op, message string) error {
	return &Error{
		Op:      op,
		Message: message,
	}
}
