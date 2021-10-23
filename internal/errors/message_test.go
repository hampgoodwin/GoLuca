package errors

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestWithMessage(t *testing.T) {
	testCases := []struct {
		description string
		err         error
		message     string
		expected    string
	}{
		{
			description: "root-error-with-message",
			err:         New("root error"),
			message:     "message error",
			expected:    "root error, with message \"message error\"",
		},
		{
			description: "nil-error-returns-nil",
			err:         nil,
		},
	}

	is := is.New(t)
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := WithMessage(tc.err, tc.message)
			if tc.err == nil {
				is.True(actual == nil)
				return
			}
			is.Equal(tc.expected, actual.Error())
			unwrapped := actual.(Message).Unwrap()
			is.Equal(tc.err, unwrapped)
			var m Message
			is.True(As(actual, &m))
			is.Equal(tc.message, m.Value)
		})
	}
}
