package errs

import "fmt"

// ErrCode represents an error in the application
type ErrCode struct {
	value int
}

// Defines a library of possible application errors
var (
	OK                 = ErrCode{value: 0}
	Cancelled          = ErrCode{value: 1}
	Unknown            = ErrCode{value: 2}
	InvalidArgument    = ErrCode{value: 3}
	DeadlineExceeded   = ErrCode{value: 4}
	NotFound           = ErrCode{value: 5}
	AlreadyExists      = ErrCode{value: 6}
	PermissionDenied   = ErrCode{value: 7}
	ResourceExhausted  = ErrCode{value: 8}
	FailedPrecondition = ErrCode{value: 9}
	Aborted            = ErrCode{value: 10}
	OutOfRange         = ErrCode{value: 11}
	Unimplemented      = ErrCode{value: 12}
	Internal           = ErrCode{value: 13}
	Unavailable        = ErrCode{value: 14}
	DataLoss           = ErrCode{value: 15}
	Unauthenticated    = ErrCode{value: 16}
)

// codeNumbers maps string representations to each ErrCode value
var codeNumbers = map[string]ErrCode{
	"ok":                  OK,
	"cancelled":           Cancelled,
	"unknown":             Unknown,
	"invalid_argument":    InvalidArgument,
	"deadline_exceeded":   DeadlineExceeded,
	"not_found":           NotFound,
	"already_exists":      AlreadyExists,
	"permission_denied":   PermissionDenied,
	"resource_exhausted":  ResourceExhausted,
	"failed_precondition": FailedPrecondition,
	"aborted":             Aborted,
	"out_of_range":        OutOfRange,
	"unimplemented":       Unimplemented,
	"internal":            Internal,
	"unavailable":         Unavailable,
	"data_loss":           DataLoss,
	"unauthenticated":     Unauthenticated,
}

// codeNames maps ErrCode values to their string representations
var codeNames [17]string

func init() {
	codeNames[OK.value] = "ok"
	codeNames[Cancelled.value] = "cancelled"
	codeNames[Unknown.value] = "unknown"
	codeNames[InvalidArgument.value] = "invalid_argument"
	codeNames[DeadlineExceeded.value] = "deadline_exceeded"
	codeNames[NotFound.value] = "not_found"
	codeNames[AlreadyExists.value] = "already_exists"
	codeNames[PermissionDenied.value] = "permission_denied"
	codeNames[ResourceExhausted.value] = "resource_exhausted"
	codeNames[FailedPrecondition.value] = "failed_precondition"
	codeNames[Aborted.value] = "aborted"
	codeNames[OutOfRange.value] = "out_of_range"
	codeNames[Unimplemented.value] = "unimplemented"
	codeNames[Internal.value] = "internal"
	codeNames[Unavailable.value] = "unavailable"
	codeNames[DataLoss.value] = "data_loss"
	codeNames[Unauthenticated.value] = "unauthenticated"
}

// Value returns the number within an ErrCode
func (ec ErrCode) Value() int {
	return ec.value
}

// String returns the text mapping of an ErrCode
func (ec ErrCode) String() string {
	return codeNames[ec.value]
}

// UnmarshalText sets the ErrCode based on its text value
func (ec *ErrCode) UnmarshalText(data []byte) error {
	errName := string(data)

	v, exists := codeNumbers[errName]
	if !exists {
		return fmt.Errorf("err code %q does not exist", errName)
	}

	*ec = v

	return nil
}

// MarshalText returns the string associated with the error code
func (ec ErrCode) MarshalText() ([]byte, error) {
	return []byte(ec.String()), nil
}

// Equal checks if two ErrCodes are equivalent
func (ec ErrCode) Equal(ec2 ErrCode) bool {
	return ec.value == ec2.value
}
