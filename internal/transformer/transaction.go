package transformer

import (
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/httpapi"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func NewTransactionFromHTTPTransaction(t httpapi.Transaction) (transaction.Transaction, error) {
	out := transaction.Transaction{}
	out.Description = t.Description
	for _, entry := range t.Entries {
		inEntry, err := NewEntryFromHTTPEntry(entry)
		if err != nil {
			return out, errors.Wrap(err, "transforming entry from httpapi transaction")
		}
		out.Entries = append(out.Entries, inEntry)
	}
	return out, nil
}

func NewEntryFromHTTPEntry(e httpapi.Entry) (transaction.Entry, error) {
	out := transaction.Entry{}
	out.Description = e.Description
	out.DebitAccount = e.DebitAccount
	out.CreditAccount = e.CreditAccount
	inAmount, err := NewAmountFromHTTPAmount(e.Amount)
	if err != nil {
		return out, errors.Wrap(err, "transforming amount from httpapi entry")
	}
	out.Amount = inAmount
	return out, nil
}
