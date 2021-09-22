package errors

import (
	"errors"
	"fmt"
)

// ErrorFlag defines a list of flags you can set on errors.
type ErrorFlag int

const (
	zero ErrorFlag = iota
	// NotFound is used when a resource is not found
	NotFound
	// NotValid indidcates a general invalid flagging
	NotValid
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
	// NoRelationshipFound indicates that a process which assumes a data relationship did not find
	// the assumed relationship
	NoRelationshipFound
)

func (ef ErrorFlag) String() string {
	return [...]string{
		"0",
		"NotFound",
		"NotValid",
		"NotValidRequest",
		"NotValidRequestData",
		"NotValidInternalData",
		"NotDeserializable",
		"NotSerializable",
		"NotRelationshipFound",
	}[ef]
}

// Flag "flags" err with an ErrorFlag that will return true from HasFlag(err, flag)
// and if the optional msg is provided, can return a msg from GetMessage(err) if set
// to exclude setting a message, provide a string zero-value of ""
func Flag(err error, flag ErrorFlag, msg string) error {
	if err == nil {
		return nil
	}
	return flagged{error: err, flag: flag, msg: msg}
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

// GetMessage returns an error message if it was set
// if no message was set, a string zero-value of "" is returned
func GetMessage(err error) string {
	for {
		if f, ok := err.(flagged); ok && f.msg != "" {
			return f.msg
		}
		if err = errors.Unwrap(err); err == nil {
			return ""
		}
	}
}

// Wrap encapsulates an error in another descriptive error allowing meaningful
// error chains.
func Wrap(err error, msg string) error {
	if f, ok := err.(flagged); ok {
		return flagged{error: fmt.Errorf("%s: %w", msg, f), flag: f.flag}
	}
	return flagged{error: fmt.Errorf("%s: %w", msg, err), flag: zero}
}

// Wrapf is a convenience function which combines print formatting using
// standard go print directives and Wrap to create more descriptive
// wrapped errors
func Wrapf(err error, msg string, a ...interface{}) error {
	return Wrap(err, fmt.Sprintf(msg, a...))
}

// FlagWrap is a convenience function which is equivalent to calling
// Wrap(Flag(error, flag), msg)
func FlagWrap(err error, flag ErrorFlag, errMsg, msg string) error {
	return Wrap(Flag(err, flag, errMsg), msg)
}

// FlagWrapf is a convenience function which is equivalent to calling
// Wrapf(Flag(err, flag), msg, a...)
func FlagWrapf(err error, flag ErrorFlag, errMsg, msg string, a ...interface{}) error {
	return Wrapf(Flag(err, flag, errMsg), msg, a...)
}

type flagged struct {
	error
	flag ErrorFlag
	msg  string
}

func (f flagged) Unwrap() error {
	return f.error
}
