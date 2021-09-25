package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
			expected:    "not valid: root error",
		},
	}

	a := assert.New(t)
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := WithError(tc.err, tc.named)
			a.Equal(tc.expected, actual.Error())
			a.True(Is(actual, tc.named))
			var w with
			a.True(As(actual, &w))
			a.Equal(tc.named, w.err)
		})
	}
}
