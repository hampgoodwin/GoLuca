package transformer

import (
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/transaction"
	httptransaction "github.com/hampgoodwin/GoLuca/pkg/http/v0/transaction"
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

func NewTransactionFromRepoTransaction(in repository.Transaction) transaction.Transaction {
	out := transaction.Transaction{}

	if in.IsZero() {
		return out
	}

	out.ID = in.ID
	out.Description = in.Description

	entries := []transaction.Entry{}
	for _, repoEntry := range in.Entries {
		entries = append(entries, NewEntryFromRepoEntry(repoEntry))
	}
	out.Entries = entries

	out.CreatedAt = in.CreatedAt

	return out
}

func NewEntryFromRepoEntry(in repository.Entry) transaction.Entry {
	out := transaction.Entry{}

	if in == (repository.Entry{}) {
		return out
	}

	out.ID = in.ID
	out.TransactionID = in.TransactionID
	out.Description = in.Description
	out.DebitAccount = in.DebitAccount
	out.CreditAccount = in.CreditAccount

	inAmount := NewAmountFromRepoAmount(in.Amount)
	out.Amount = inAmount

	out.CreatedAt = in.CreatedAt

	return out
}

func NewRepoTransactionFromTransaction(in transaction.Transaction) repository.Transaction {
	out := repository.Transaction{}

	if in.IsZero() {
		return out
	}

	out.ID = in.ID
	out.Description = in.Description

	entries := []repository.Entry{}
	for _, entry := range in.Entries {
		entries = append(entries, NewRepoEntryFromEntry(entry))
	}
	out.Entries = entries

	out.CreatedAt = in.CreatedAt

	return out
}

func NewRepoEntryFromEntry(in transaction.Entry) repository.Entry {
	out := repository.Entry{}

	if in == (transaction.Entry{}) {
		return out
	}

	out.ID = in.ID
	out.TransactionID = in.TransactionID
	out.Description = in.Description
	out.DebitAccount = in.DebitAccount
	out.CreditAccount = in.CreditAccount

	inAmount := NewRepoAmountFromAmount(in.Amount)
	out.Amount = inAmount

	out.CreatedAt = in.CreatedAt

	return out
}
