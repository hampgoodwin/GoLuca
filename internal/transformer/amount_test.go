package transformer

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/go-playground/validator"
	"github.com/hampgoodwin/GoLuca/pkg/amount"
	httpamount "github.com/hampgoodwin/GoLuca/pkg/http/v0/amount"
	"github.com/hampgoodwin/errors"
	"github.com/matryer/is"
)

func TestNewAmountFromHTTPAmount(t *testing.T) {
	testCases := []struct {
		description string
		httpamount  httpamount.Amount
		expected    amount.Amount
		err         error
	}{
		{description: "empty"},
		{
			description: "value-overflows-int64",
			httpamount:  httpamount.Amount{Value: "9223372036854775808"},
			err:         &strconv.NumError{},
		},
		{
			description: "invalid-value",
			httpamount:  httpamount.Amount{Value: "-10"},
			err:         validator.ValidationErrors{},
		},
		{
			description: "invalid-currency",
			httpamount:  httpamount.Amount{Currency: "KRKZ"},
			err:         validator.ValidationErrors{},
		},
		{
			description: "success",
			httpamount:  httpamount.Amount{Value: "100", Currency: "USD"},
			expected:    amount.Amount{Value: 100, Currency: "USD"},
		},
	}

	a := is.New(t)

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual, err := NewAmountFromHTTPAmount(tc.httpamount)
			if err != nil {
				a.True(errors.As(err, &tc.err))
				return
			}
			a.NoErr(err)
			a.Equal(tc.expected, actual)
		})
	}
}
