package transformer

import (
	"fmt"

	modelv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/model/v1"
	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/amount"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/transaction"
	httptransaction "github.com/hampgoodwin/GoLuca/pkg/http/v0/transaction"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			return out, fmt.Errorf("transforming entry from http create transaction: %w", err)
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
		return out, fmt.Errorf("transforming amount from http create entry: %w", err)
	}
	out.Amount = inAmount

	return out, nil
}

func NewHTTPTransactionFromTransaction(in transaction.Transaction) httptransaction.Transaction {
	out := httptransaction.Transaction{}

	if in.IsZero() {
		return out
	}

	out.ID = in.ID
	out.Description = in.Description

	for _, entry := range in.Entries {
		inEntry := NewHTTPEntryFromEntry(entry)
		out.Entries = append(out.Entries, inEntry)
	}

	out.CreatedAt = in.CreatedAt

	return out
}

func NewHTTPEntryFromEntry(in transaction.Entry) httptransaction.Entry {
	out := httptransaction.Entry{}

	if in == (transaction.Entry{}) {
		return out
	}

	out.ID = in.ID
	out.TransactionID = in.TransactionID
	out.Description = in.Description
	out.DebitAccount = in.DebitAccount
	out.CreditAccount = in.CreditAccount

	inAmount := NewHTTPAmountFromAmount(in.Amount)
	out.Amount = inAmount

	out.CreatedAt = in.CreatedAt

	return out
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

func NewProtoTransactionFromTransaction(in transaction.Transaction) *modelv1.Transaction {
	if in.IsZero() {
		return nil
	}

	out := &modelv1.Transaction{
		Id:          in.ID,
		Description: in.Description,
		Entries:     nil, // filled in after
		CreatedAt:   timestamppb.New(in.CreatedAt),
	}

	var entries []*modelv1.Entry
	for _, entry := range in.Entries {
		entry := NewProtoEntryFromEntry(entry)
		entries = append(entries, entry)
	}
	if len(entries) > 0 {
		out.Entries = entries
	}

	return out
}

func NewProtoEntryFromEntry(in transaction.Entry) *modelv1.Entry {
	if in == (transaction.Entry{}) {
		return nil
	}

	out := &modelv1.Entry{
		Id:            in.ID,
		TransactionId: in.TransactionID,
		Description:   in.Description,
		DebitAccount:  in.DebitAccount,
		CreditAccount: in.CreditAccount,
		Amount:        NewProtoAmountFromAmount(in.Amount),
		CreatedAt:     timestamppb.New(in.CreatedAt),
	}

	return out
}

func NewProtoAmountFromAmount(in amount.Amount) *modelv1.Amount {
	if in == (amount.Amount{}) {
		return nil
	}

	out := &modelv1.Amount{
		Value:    in.Value,
		Currency: in.Currency,
	}

	return out
}

func NewTransactionFromProtoCreateTransaction(in *servicev1.CreateTransactionRequest) transaction.Transaction {
	out := transaction.Transaction{}

	if in == nil {
		return out
	}

	out.Description = in.GetDescription()

	for _, entry := range in.Entries {
		out.Entries = append(out.Entries, NewEntryFromProtoCreateEntry(entry))
	}

	return out
}

func NewEntryFromProtoCreateEntry(in *servicev1.CreateEntry) transaction.Entry {
	out := transaction.Entry{}

	if in == nil {
		return out
	}

	out.Description = in.GetDescription()
	out.DebitAccount = in.GetDebitAccount()
	out.CreditAccount = in.GetCreditAccount()
	out.Amount = NewAmountFromProtoAmount(in.Amount)

	return out
}
