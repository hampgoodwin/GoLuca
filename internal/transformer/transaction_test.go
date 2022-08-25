package transformer

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/pkg/amount"
	httpamount "github.com/hampgoodwin/GoLuca/pkg/http/amount"
	httptransaction "github.com/hampgoodwin/GoLuca/pkg/http/transaction"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/hampgoodwin/errors"
	"github.com/matryer/is"
)

func TestNewTransactionFromHTTPCreateTransaction(t *testing.T) {
	debitAccount := uuid.NewString()
	creditAccount := uuid.NewString()
	testCases := []struct {
		description     string
		httpTransaction httptransaction.CreateTransaction
		expected        transaction.Transaction
		err             error
	}{
		{
			description: "empty",
		},
		{
			description: "success",
			httpTransaction: httptransaction.CreateTransaction{
				Description: "transaction",
				Entries: []httptransaction.CreateEntry{
					{
						Description:   "1",
						DebitAccount:  debitAccount,
						CreditAccount: creditAccount,
						Amount: httpamount.Amount{
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
			actual, err := NewTransactionFromHTTPCreateTransaction(tc.httpTransaction)
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
