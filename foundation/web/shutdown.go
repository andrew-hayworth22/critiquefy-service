package web

import "errors"

// shutdownError represents an error that should shut the system down
type shutdownError struct {
	Message string
}

// NewShutdownError constructs a new shutdownError
func NewShutdownError(message string) error {
	return &shutdownError{message}
}

// Error allows shutdownErrors to satisfy the error interface
func (se *shutdownError) Error() string {
	return se.Message
}

// IsShutdown checks if an error should shut the system down
func IsShutdown(err error) bool {
	var se *shutdownError
	return errors.As(err, &se)
}
