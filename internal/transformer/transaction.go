package transformer

import (
	"github.com/hampgoodwin/GoLuca/internal/http/api"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/hampgoodwin/errors"
)

func NewTransactionFromHTTPTransaction(t api.Transaction) (transaction.Transaction, error) {
	out := transaction.Transaction{}
	if t.IsZero() {
		return out, nil
	}
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

func NewEntryFromHTTPEntry(e api.Entry) (transaction.Entry, error) {
	out := transaction.Entry{}
	if e == (api.Entry{}) {
		return out, nil
	}
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
