package errors

import (
	"fmt"
)

// WithError adds a custom error alongside the original error which can
// later be reported as found in the error chain via .Is()
func WithError(rootErr, withErr error) error {
	return with{error: rootErr, err: withErr}
}

type with struct {
	error
	err error
}

func (w with) Error() string { return fmt.Sprintf("%s, with error %q", w.error.Error(), w.err.Error()) }

func (w with) Unwrap() error {
	return w.error
}

func (w with) Is(err error) bool {
	return w.error == err || w.err == err
}
