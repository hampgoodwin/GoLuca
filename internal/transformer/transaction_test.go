package transformer

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hampgoodwin/GoLuca/internal/amount"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/transaction"
	httpamount "github.com/hampgoodwin/GoLuca/pkg/http/v0/amount"
	httptransaction "github.com/hampgoodwin/GoLuca/pkg/http/v0/transaction"

	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestNewTransactionFromHTTPCreateTransaction(t *testing.T) {
	debitUUIDV7, _ := uuid.NewV7()
	debitAccount := debitUUIDV7.String()
	creditUUIDV7, _ := uuid.NewV7()
	creditAccount := creditUUIDV7.String()
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

func TestNewHTTPTransactionFromTransaction(t *testing.T) {
	transactionUUIDV7, _ := uuid.NewV7()
	transactionID := transactionUUIDV7.String()

	entryUUIDV7, _ := uuid.NewV7()
	entryID := entryUUIDV7.String()

	debitUUIDV7, _ := uuid.NewV7()
	debitAccount := debitUUIDV7.String()
	creditUUIDV7, _ := uuid.NewV7()
	creditAccount := creditUUIDV7.String()
	testCases := []struct {
		description string
		transaction transaction.Transaction
		expected    httptransaction.Transaction
		err         error
	}{
		{
			description: "empty",
		},
		{
			description: "success",
			transaction: transaction.Transaction{
				ID:          transactionID,
				Description: "transaction",
				Entries: []transaction.Entry{
					{
						ID:            entryID,
						TransactionID: transactionID,
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
			expected: httptransaction.Transaction{
				ID:          transactionID,
				Description: "transaction",
				Entries: []httptransaction.Entry{
					{
						ID:            entryID,
						TransactionID: transactionID,
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
		},
	}

	a := is.New(t)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := NewHTTPTransactionFromTransaction(tc.transaction)
			a.Equal(tc.expected, actual)
		})
	}
}

func TestNewTransactionFromRepoTransaction(t *testing.T) {
	transactionUUIDV7, _ := uuid.NewV7()
	transactionID := transactionUUIDV7.String()

	entryUUIDV7, _ := uuid.NewV7()
	entryID := entryUUIDV7.String()

	debitUUIDV7, _ := uuid.NewV7()
	debitAccount := debitUUIDV7.String()
	creditUUIDV7, _ := uuid.NewV7()
	creditAccount := creditUUIDV7.String()
	testCases := []struct {
		description string
		transaction repository.Transaction
		expected    transaction.Transaction
		err         error
	}{
		{
			description: "empty",
		},
		{
			description: "success",
			transaction: repository.Transaction{
				ID:          transactionID,
				Description: "transaction",
				Entries: []repository.Entry{
					{
						ID:            entryID,
						TransactionID: transactionID,
						Description:   "1",
						DebitAccount:  debitAccount,
						CreditAccount: creditAccount,
						Amount: repository.Amount{
							Value:    100,
							Currency: "USD",
						},
					},
				},
			},
			expected: transaction.Transaction{
				ID:          transactionID,
				Description: "transaction",
				Entries: []transaction.Entry{
					{
						ID:            entryID,
						TransactionID: transactionID,
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
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := NewTransactionFromRepoTransaction(tc.transaction)
			a.Equal(tc.expected, actual)
		})
	}
}

func TestNewRepoTransactionFromTransaction(t *testing.T) {
	transactionUUIDV7, _ := uuid.NewV7()
	transactionID := transactionUUIDV7.String()

	entryUUIDV7, _ := uuid.NewV7()
	entryID := entryUUIDV7.String()

	debitUUIDV7, _ := uuid.NewV7()
	debitAccount := debitUUIDV7.String()
	creditUUIDV7, _ := uuid.NewV7()
	creditAccount := creditUUIDV7.String()
	testCases := []struct {
		description string
		transaction transaction.Transaction
		expected    repository.Transaction
		err         error
	}{
		{
			description: "empty",
		},
		{
			description: "success",
			transaction: transaction.Transaction{
				ID:          transactionID,
				Description: "transaction",
				Entries: []transaction.Entry{
					{
						ID:            entryID,
						TransactionID: transactionID,
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
			expected: repository.Transaction{
				ID:          transactionID,
				Description: "transaction",
				Entries: []repository.Entry{
					{
						ID:            entryID,
						TransactionID: transactionID,
						Description:   "1",
						DebitAccount:  debitAccount,
						CreditAccount: creditAccount,
						Amount: repository.Amount{
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
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := NewRepoTransactionFromTransaction(tc.transaction)
			a.Equal(tc.expected, actual)
		})
	}
}
