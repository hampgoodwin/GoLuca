package transformer

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hampgoodwin/GoLuca/internal/amount"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	httpamount "github.com/hampgoodwin/GoLuca/pkg/http/v0/amount"

	"github.com/go-playground/validator/v10"
	"github.com/matryer/is"
)

func TestNewAmountFromHTTPAmount(t *testing.T) {
	testCases := []struct {
		description string
		httpamount  httpamount.Amount
		expected    amount.Amount
		assertErr   func(err error) bool
		err         error
	}{
		{description: "empty"},
		{
			description: "value-overflows-int64",
			httpamount:  httpamount.Amount{Value: "9223372036854775808"},
			assertErr: func(err error) bool {
				var numErr *strconv.NumError
				return errors.As(err, &numErr)
			},
		},
		{
			description: "invalid-value",
			httpamount:  httpamount.Amount{Value: "-10"},
			assertErr: func(err error) bool {
				var validationErrors validator.ValidationErrors
				return errors.As(err, &validationErrors)
			},
		},
		{
			description: "invalid-currency",
			httpamount:  httpamount.Amount{Value: "0", Currency: "KRKZ"},
			assertErr: func(err error) bool {
				var validationErrors validator.ValidationErrors
				b := errors.As(err, &validationErrors)
				return b
			},
		},
		{
			description: "success",
			httpamount:  httpamount.Amount{Value: "100", Currency: "USD"},
			expected:    amount.Amount{Value: 100, Currency: "USD"},
		},
	}

	a := is.New(t)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual, err := NewAmountFromHTTPAmount(tc.httpamount)
			if tc.assertErr != nil {
				a.True(tc.assertErr(err))
				return
			}
			a.NoErr(err)
			a.Equal(tc.expected, actual)
		})
	}
}

func TestNewHTTPAmountFromAmount(t *testing.T) {
	testCases := []struct {
		description string
		amount      amount.Amount
		expected    httpamount.Amount
		err         error
	}{
		{description: "empty"},
		{
			description: "success",
			amount:      amount.Amount{Value: 100, Currency: "USD"},
			expected:    httpamount.Amount{Value: "100", Currency: "USD"},
		},
	}

	a := is.New(t)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := NewHTTPAmountFromAmount(tc.amount)
			a.Equal(tc.expected, actual)
		})
	}
}

func TestNewAmountFromRepoAmount(t *testing.T) {
	testCases := []struct {
		description string
		amount      repository.Amount
		expected    amount.Amount
		err         error
	}{
		{description: "empty"},
		{
			description: "success",
			amount:      repository.Amount{Value: 100, Currency: "USD"},
			expected:    amount.Amount{Value: 100, Currency: "USD"},
		},
	}

	a := is.New(t)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := NewAmountFromRepoAmount(tc.amount)
			a.Equal(tc.expected, actual)
		})
	}
}

func TestNewRepoAmountFromAmount(t *testing.T) {
	testCases := []struct {
		description string
		amount      amount.Amount
		expected    repository.Amount
		err         error
	}{
		{description: "empty"},
		{
			description: "success",
			amount:      amount.Amount{Value: 100, Currency: "USD"},
			expected:    repository.Amount{Value: 100, Currency: "USD"},
		},
	}

	a := is.New(t)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := NewRepoAmountFromAmount(tc.amount)
			a.Equal(tc.expected, actual)
		})
	}
}
