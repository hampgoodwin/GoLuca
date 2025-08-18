package errors

import "fmt"

type NotFoundErr struct {
	ID   string
	Type string
}

func (e NotFoundErr) Error() string {
	return fmt.Sprintf("%q %q not found", e.Type, e.ID)
}
