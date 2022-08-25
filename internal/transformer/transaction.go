package transformer

import (
	httptransaction "github.com/hampgoodwin/GoLuca/pkg/http/v0/transaction"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/hampgoodwin/errors"
)

func NewTransactionFromHTTPCreateTransaction(in httptransaction.CreateTransaction) (transaction.Transaction, error) {
	out := transaction.Transaction{}
	if in.IsZero() {
		return out, nil
	}
	out.Description = in.Description
	for _, entry := range in.Entries {
		inEntry, err := NewEntryFromHTTPCreateEntry(entry)
		if err != nil {
			return out, errors.Wrap(err, "transforming entry from http create transaction")
		}
		out.Entries = append(out.Entries, inEntry)
	}
	return out, nil
}

func NewEntryFromHTTPCreateEntry(in httptransaction.CreateEntry) (transaction.Entry, error) {
	out := transaction.Entry{}
	if in == (httptransaction.CreateEntry{}) {
		return out, nil
	}
	out.Description = in.Description
	out.DebitAccount = in.DebitAccount
	out.CreditAccount = in.CreditAccount
	inAmount, err := NewAmountFromHTTPAmount(in.Amount)
	if err != nil {
		return out, errors.Wrap(err, "transforming amount from http create entry")
	}
	out.Amount = inAmount
	return out, nil
}

func NewHTTPTransactionFromTransaction(in transaction.Transaction) httptransaction.Transaction {
	return httptransaction.Transaction(in)
}
