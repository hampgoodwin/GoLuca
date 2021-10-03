package errors

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestWithError(t *testing.T) {
	testCases := []struct {
		description string
		err         error
		named       error
		expected    string
	}{
		{
			description: "root-error-with-named-error",
			err:         New("root error"),
			named:       NotValid,
			expected:    "root error, with error \"not valid\"",
		},
	}

	is := is.New(t)
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := WithError(tc.err, tc.named)
			is.Equal(tc.expected, actual.Error())
			is.True(Is(actual, tc.named))
			var w with
			is.True(As(actual, &w))
			is.Equal(tc.named, w.err)
		})
	}
}
