package transformer

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/httpapi"
	"github.com/hampgoodwin/GoLuca/pkg/amount"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/stretchr/testify/require"
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

	a := require.New(t)
	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual, err := NewTransactionFromHTTPTransaction(tc.httpTransaction)
			if tc.err != nil {
				a.Error(err)
				a.True(errors.Is(err, tc.err))
				return
			}
			a.NoError(err)
			a.Equal(tc.expected, actual)
		})
	}
}
