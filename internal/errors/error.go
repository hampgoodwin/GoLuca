package errors

import (
	"fmt"
)

var (
	// NotFound is used when a resource is not found
	NotFound = New("not found")
	// NotValid indidcates a general invalid error
	NotValid = New("not valid")
	// NotValidRequest indicates that something other than a requests body is invalid
	// for example, if a request is on http protocol, maybe a header, or query parameter
	// is invalid
	NotValidRequest = New("not valid request")
	// NotValidRequestData indicates that a request is valid, but the data
	// provided in the request is invalid
	NotValidRequestData = New("not valid request data")
	// NotValidInternalData indicates a successful request, but that the application
	// data is malformed
	NotValidInternalData = New("not valid internal data")
	// NotDeserializable indicates provided or internal data was not successfully deserialized to
	// application data structures
	NotDeserializable = New("not deserializable")
	// NotSerializable indicates provided or internals data was not successfully serialized to
	// application interface data structures
	NotSerializable = New("not serializable")
	// NoRelationshipFound indicates that a process which assumes a data relationship did not find
	// the assumed relationship
	NoRelationshipFound = New("not relationship found")
)

func WithError(err1, err2 error) error {
	return with{error: err1, err: err2}
}

type with struct {
	error
	err error
}

func (w with) Error() string { return fmt.Sprintf("%s: %s", w.err.Error(), w.error.Error()) }

func (w *with) Unwrap() error { return w.error }

func (w with) Is(err error) bool {
	return w.error == err || w.err == err
}
