package service

import (
	"github.com/hampgoodwin/GoLuca/internal/transaction"
	"github.com/hampgoodwin/GoLuca/internal/transaction/repository"
)

func newRepoTransactionFromTransaction(in transaction.Transaction) repository.Transaction {
	out := repository.Transaction{}

	if in.IsZero() {
		return out
	}

	out.ID = in.ID
	out.Description = in.Description

	entries := []repository.Entry{}
	for _, entry := range in.Entries {
		entries = append(entries, newRepoEntryFromEntry(entry))
	}
	out.Entries = entries

	out.CreatedAt = in.CreatedAt

	return out
}

func newRepoEntryFromEntry(in transaction.Entry) repository.Entry {
	out := repository.Entry{}

	if in == (transaction.Entry{}) {
		return out
	}

	out.ID = in.ID
	out.TransactionID = in.TransactionID
	out.Description = in.Description
	out.DebitAccount = in.DebitAccount
	out.CreditAccount = in.CreditAccount

	inAmount := newRepoAmountFromAmount(in.Amount)
	out.Amount = inAmount

	out.CreatedAt = in.CreatedAt

	return out
}

func newTransactionFromRepoTransaction(in repository.Transaction) transaction.Transaction {
	out := transaction.Transaction{}

	if in.IsZero() {
		return out
	}

	out.ID = in.ID
	out.Description = in.Description

	entries := []transaction.Entry{}
	for _, repoEntry := range in.Entries {
		entries = append(entries, newEntryFromRepoEntry(repoEntry))
	}
	out.Entries = entries

	out.CreatedAt = in.CreatedAt

	return out
}

func newEntryFromRepoEntry(in repository.Entry) transaction.Entry {
	out := transaction.Entry{}

	if in == (repository.Entry{}) {
		return out
	}

	out.ID = in.ID
	out.TransactionID = in.TransactionID
	out.Description = in.Description
	out.DebitAccount = in.DebitAccount
	out.CreditAccount = in.CreditAccount

	inAmount := newAmountFromRepoAmount(in.Amount)
	out.Amount = inAmount

	out.CreatedAt = in.CreatedAt

	return out
}

func newAmountFromRepoAmount(in repository.Amount) transaction.Amount {
	out := transaction.Amount{}

	if in == (repository.Amount{}) {
		return out
	}

	out.Value = in.Value
	out.Currency = in.Currency

	return out
}

func newRepoAmountFromAmount(in transaction.Amount) repository.Amount {
	out := repository.Amount{}

	if in == (transaction.Amount{}) {
		return out
	}

	out.Value = in.Value
	out.Currency = in.Currency

	return out
}
