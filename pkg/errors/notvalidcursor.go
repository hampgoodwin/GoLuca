package errors

import "fmt"

type NotValidCursorErr struct {
	Cursor string
}

func (e NotValidCursorErr) Error() string {
	return fmt.Sprintf("invalid cursor %q", e.Cursor)
}
