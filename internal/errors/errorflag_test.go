package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringer(t *testing.T) {
	testCases := []struct {
		input    ErrorFlag
		expected string
	}{
		{zero, "0"},
		{NotFound, "NotFound"},
		{NotValidRequest, "NotValidRequest"},
		{NotValidRequestData, "NotValidRequestData"},
		{NotValidInternalData, "NotValidInternalData"},
	}

	a := assert.New(t)

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.input), func(t *testing.T) {
			t.Parallel()
			actual := tc.input.String()
			a.Equal(tc.expected, actual)
		})
	}
}
