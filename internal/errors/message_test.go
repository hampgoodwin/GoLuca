package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

	a := assert.New(t)
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := WithMessage(tc.err, tc.message)
			if tc.err == nil {
				a.Nil(actual)
				return
			}
			a.Equal(tc.expected, actual.Error())
			unwrapped := actual.(Message).Unwrap()
			a.Equal(tc.err, unwrapped)
			var m Message
			a.True(As(actual, &m))
			a.Equal(tc.message, m.Value)
		})
	}
}
