package errors

import "fmt"

// Wrap applies a consistent formatting of `fmt.Errorf("%s: %w")`
// with the provided values
func Wrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
