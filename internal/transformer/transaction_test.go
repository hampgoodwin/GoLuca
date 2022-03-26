package transformer

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/internal/httpapi"
	"github.com/hampgoodwin/GoLuca/pkg/amount"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/hampgoodwin/errors"
	"github.com/matryer/is"
)

func TestNewTransactionFromHTTPTransaction(t *testing.T) {
	debitAccount := uuid.NewString()
	creditAccount := uuid.NewString()
	testCases := []struct {
		description     string
		httpTransaction httpapi.Transaction
		expected        transaction.Transaction
		err             error
	}{
		{
			description: "empty",
		},
		{
			description: "success",
			httpTransaction: httpapi.Transaction{
				Description: "transaction",
				Entries: []httpapi.Entry{
					{
						Description:   "1",
						DebitAccount:  debitAccount,
						CreditAccount: creditAccount,
						Amount: httpapi.Amount{
							Value:    "100",
							Currency: "USD",
						},
					},
				},
			},
			expected: transaction.Transaction{
				Description: "transaction",
				Entries: []transaction.Entry{
					{
						Description:   "1",
						DebitAccount:  debitAccount,
						CreditAccount: creditAccount,
						Amount: amount.Amount{
							Value:    100,
							Currency: "USD",
						},
					},
				},
			},
		},
	}

	a := is.New(t)
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual, err := NewTransactionFromHTTPTransaction(tc.httpTransaction)
			if tc.err != nil {
				a.True(err != nil)
				a.True(errors.Is(err, tc.err))
				return
			}
			a.NoErr(err)
			a.Equal(tc.expected, actual)
		})
	}
}
