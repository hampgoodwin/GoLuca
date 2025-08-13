package errors

import "errors"

// Sentinel Errors
var (
	// ErrNotFound is used when a resource is not found
	ErrNotFound = errors.New("not found")
	// ErrNotValid indidcates a general invalid error
	ErrNotValid = errors.New("not valid")
	// ErrNotValidRequest indicates that something other than a requests body is invalid
	// for example, if a request is on http protocol, maybe a header, or query parameter
	// is invalid
	ErrNotValidRequest = errors.New("not valid request")
	// ErrNotValidRequestData indicates that a request is valid, but the data
	// provided in the request is invalid
	ErrNotValidRequestData = errors.New("not valid request data")
	// ErrNotValidInternalData indicates a successful request, but that the application
	// data is malformed
	ErrNotValidInternalData = errors.New("not valid internal data")
	// ErrNotDeserializable indicates provided or internal data was not successfully deserialized to
	// application data structures
	ErrNotDeserializable = errors.New("not deserializable")
	// ErrNotSerializable indicates provided or internals data was not successfully serialized to
	// application interface data structures
	ErrNotSerializable = errors.New("not serializable")
	// NoRelationshipFound indicates that a process which assumes a data relationship did not find
	// the assumed relationship
	ErrNoRelationshipFound = errors.New("not relationship found")
	// ErrNotKnown indicates an application failure for which the failure is not known
	ErrNotKnown = errors.New("not known")
	// ErrNotUnique indicates an expectations of uniqueness which is not met
	ErrNotUnique = errors.New("not unique")
)
