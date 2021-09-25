package errors

import "fmt"

type Message struct {
	error
	Value string
}

func (m Message) Error() string { return fmt.Sprintf("%s, with message %q", m.error.Error(), m.Value) }

func (m Message) Unwrap() error { return m.error }

func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return Message{error: err, Value: message}
}
