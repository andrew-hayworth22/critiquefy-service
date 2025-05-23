package errs

import (
	"errors"
	"fmt"
)

// Error represents an application level error
type Error struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

// New creates a new application Error structure with an existing error
func New(code ErrCode, err error) Error {
	return Error{
		Code:    code,
		Message: err.Error(),
	}
}

// Newf creates a new application Error structure with the ability to format the message
func Newf(code ErrCode, format string, v ...any) Error {
	return Error{
		Code:    code,
		Message: fmt.Sprintf(format, v...),
	}
}

// Error returns the error message and satisfies the error interface
func (err Error) Error() string {
	return err.Message
}

// IsError checks if a given error is an application Error
func IsError(err error) bool {
	var er Error
	return errors.As(err, &er)
}

// GetError converts an error into an application Error
func GetError(err error) Error {
	var er Error
	if !errors.As(err, &er) {
		return Error{}
	}
	return er
}
