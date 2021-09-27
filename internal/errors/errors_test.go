package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithErrorWithMessage(t *testing.T) {
	a := require.New(t)
	rootErr := New("root error")
	withError := WithError(rootErr, NotValid)

	a.Equal("root error, with error \"not valid\"", withError.Error())
	a.True(Is(withError, NotValid))

	withMessage := WithMessage(withError, "error message")
	a.Equal("root error, with error \"not valid\"", withError.Error())
	a.True(Is(withError, NotValid))
	var m Message
	a.True(As(withMessage, &m))
	a.Equal("error message", m.Value)

	wrap1 := Wrap(withMessage, "wrap1")
	a.Equal("wrap1: root error, with error \"not valid\", with message \"error message\"", wrap1.Error())
}

type customError struct {
	Val string
}

func (ce customError) Error() string {
	return ce.Val
}

func TestWithErrorWithMessage_CustomRootError(t *testing.T) {
	a := require.New(t)
	rootErr := customError{Val: "custom error"}
	withError := WithError(rootErr, NotValid)

	a.Equal("custom error, with error \"not valid\"", withError.Error())
	a.True(Is(withError, NotValid))
	var ce customError
	a.True(As(withError, &ce))

	withMessage := WithMessage(withError, "error message")
	a.Equal("custom error, with error \"not valid\"", withError.Error())
	a.True(Is(withError, NotValid))
	var m Message
	a.True(As(withMessage, &m))
	a.Equal("error message", m.Value)

	wrap1 := Wrap(withMessage, "wrap1")
	a.Equal("wrap1: custom error, with error \"not valid\", with message \"error message\"", wrap1.Error())
}

func TestWithErrorMessage(t *testing.T) {
	a := require.New(t)
	rootErr := New("root error")
	withErrorMessage := WithErrorMessage(rootErr, NotValid, "message")

	a.Equal("root error, with error \"not valid\", with message \"message\"",
		withErrorMessage.Error())
}
func TestWithErrorMessage_CustomRootError(t *testing.T) {
	a := require.New(t)
	rootErr := customError{Val: "custom error"}
	withErrorMessage := WithErrorMessage(rootErr, NotValid, "message")

	a.Equal("custom error, with error \"not valid\", with message \"message\"",
		withErrorMessage.Error())
}

func TestWrapWithErrorMessage(t *testing.T) {
	a := require.New(t)
	rootErr := New("root error")
	errorMsg := "message"
	wrapMsg := "wrap"
	withErrorMessage := WrapWithErrorMessage(rootErr, NotValid, errorMsg, wrapMsg)

	a.Equal("wrap: root error, with error \"not valid\", with message \"message\"",
		withErrorMessage.Error())
}
