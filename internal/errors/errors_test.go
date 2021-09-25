package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithErrorWithMessage(t *testing.T) {
	a := require.New(t)
	rootErr := New("root error")
	withError := WithError(rootErr, NotValid)

	a.Equal("not valid: root error", withError.Error())
	a.True(Is(withError, NotValid))

	withMessage := WithMessage(withError, "error message")
	a.Equal("not valid: root error", withError.Error())
	a.True(Is(withError, NotValid))
	var m Message
	a.True(As(withMessage, &m))
	a.Equal("error message", m.Value)

	wrap1 := Wrap(withMessage, "wrap1")
	a.Equal("wrap1: not valid: root error, with message \"error message\"", wrap1.Error())
}
