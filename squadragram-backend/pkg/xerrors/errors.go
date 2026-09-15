package xerrors

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	errValidation = errors.New("validation error")
	errBusiness   = errors.New("business error")
)
var errorcodeCounter int64 = 1000

type SentinelError struct {
	Msg  string
	Code int64
}

func NewSentinel(msg string) *SentinelError {
	newcode := atomic.AddInt64(&errorcodeCounter, 1)
	return &SentinelError{
		Msg:  msg,
		Code: newcode,
	}
}

func NewValidationFromError(err error) error {
	return fmt.Errorf("%w: %w", errValidation, err)
}

func NewValidationFromBusiness(err error) error {
	return fmt.Errorf("%w: %w", errBusiness, err)
}

func IsValidation(err error) bool {
	return errors.Is(err, errValidation)
}

func IsBusiness(err error) bool {
	return errors.Is(err, errBusiness)
}
