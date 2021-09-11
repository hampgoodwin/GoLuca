package errors

import (
	"errors"
	"fmt"
)

// ErrorFlag defines a list of flags you can set on errors.
type ErrorFlag int

const (
	zero = iota
	// NotFound is used when a resource is not found
	NotFound
	// NotValidRequest indicates that something other than a requests body is invalid
	// for example, if a request is on http protocol, maybe a header, or query parameter
	// is invalid
	NotValidRequest
	// NotValidRequestData indicates that a request is valid, but the data
	// provided in the request is invalid
	NotValidRequestData
	// NotValidInternalData indicates a successful request, but that the application
	// data is malformed
	NotValidInternalData
	// NotDeserializable indicates provided or internal data was not successfully deserialized to
	// application data structures
	NotDeserializable
	// NotSerializable indicates provided or internals data was not successfully serialized to
	// application interface data structures
	NotSerializable
)

func (ef ErrorFlag) String() string {
	return [...]string{
		"0",
		"NotFound",
		"NotValidRequest",
		"NotValidRequestData",
		"NotValidInternalData",
	}[ef]
}

// Flag wraps err with an error that will return true from HasFlag(err, flag).
func Flag(err error, flag ErrorFlag) error {
	if err == nil {
		return nil
	}
	return flagged{error: err, flag: flag}
}

// HasFlag reports if err has been flagged with the given flag.
func HasFlag(err error, flag ErrorFlag) bool {
	for {
		if f, ok := err.(flagged); ok && f.flag == flag {
			return true
		}
		if err = errors.Unwrap(err); err == nil {
			return false
		}
	}
}

func Wrap(err error, msg string) error {
	if f, ok := err.(flagged); ok {
		return flagged{error: fmt.Errorf("%s: %w", msg, f), flag: f.flag}
	}
	return flagged{error: fmt.Errorf("%s: %w", msg, err), flag: zero}
}

func Wrapf(err error, msg string, a ...interface{}) error {
	if f, ok := err.(flagged); ok {
		return flagged{error: fmt.Errorf("%s: %w", fmt.Sprintf(msg, a...), f), flag: f.flag}
	}
	return flagged{error: fmt.Errorf("%s: %w", fmt.Sprintf(msg, a...), err), flag: zero}
}

func WrapFlag(err error, msg string, flag ErrorFlag) error {
	return Flag(fmt.Errorf("%s: %w", msg, err), flag)
}

func WrapfFlag(err error, msg string, flag ErrorFlag, a ...interface{}) error {
	return Flag(fmt.Errorf("%s: %w", fmt.Sprintf(msg, a...), err), flag)
}

type flagged struct {
	error
	flag ErrorFlag
}

func (f flagged) Unwrap() error {
	return f.error
}
