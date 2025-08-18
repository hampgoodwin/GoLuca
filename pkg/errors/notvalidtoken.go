package errors

import "fmt"

type NotValidTokenErr struct {
	Token string
}

func (e NotValidTokenErr) Error() string {
	return fmt.Sprintf("invalid token %q", e.Token)
}
