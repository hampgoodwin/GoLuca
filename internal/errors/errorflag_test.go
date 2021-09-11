package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringer(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		input    ErrorFlag
		expected string
	}{
		{zero, "0"},
		{NotFound, "NotFound"},
		{NotValid, "NotValid"},
		{NotValidRequest, "NotValidRequest"},
		{NotValidRequestData, "NotValidRequestData"},
		{NotValidInternalData, "NotValidInternalData"},
		{NotDeserializable, "NotDeserializable"},
		{NotSerializable, "NotSerializable"},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.input), func(t *testing.T) {
			t.Parallel()
			actual := tc.input.String()
			a.Equal(tc.expected, actual)
		})
	}
}
