package transaction

import (
	entryv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/entry/v1"
	transactionv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/transaction/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewProtoTransactionFromTransaction(in Transaction) *transactionv1.Transaction {
	if in.IsZero() {
		return nil
	}

	out := &transactionv1.Transaction{
		Id:          in.ID,
		Description: in.Description,
		Entries:     nil, // filled in after
		CreatedAt:   timestamppb.New(in.CreatedAt),
	}

	var entries []*entryv1.Entry
	for _, entry := range in.Entries {
		entry := NewProtoEntryFromEntry(entry)
		entries = append(entries, entry)
	}
	if len(entries) > 0 {
		out.Entries = entries
	}

	return out
}

func NewProtoEntryFromEntry(in Entry) *entryv1.Entry {
	if in == (Entry{}) {
		return nil
	}

	out := &entryv1.Entry{
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

func NewProtoAmountFromAmount(in Amount) *entryv1.Amount {
	if in == (Amount{}) {
		return nil
	}

	out := &entryv1.Amount{
		Value:    in.Value,
		Currency: in.Currency,
	}

	return out
}

func NewTransactionFromProtoCreateTransaction(in *transactionv1.CreateTransactionRequest) Transaction {
	out := Transaction{}

	if in == nil {
		return out
	}

	out.Description = in.GetDescription()

	for _, entry := range in.Entries {
		out.Entries = append(out.Entries, NewEntryFromProtoCreateEntry(entry))
	}

	return out
}

func NewEntryFromProtoCreateEntry(in *transactionv1.CreateEntry) Entry {
	out := Entry{}

	if in == nil {
		return out
	}

	out.Description = in.GetDescription()
	out.DebitAccount = in.GetDebitAccount()
	out.CreditAccount = in.GetCreditAccount()
	out.Amount = NewAmountFromProtoAmount(in.Amount)

	return out
}

func NewAmountFromProtoAmount(in *entryv1.Amount) Amount {
	out := Amount{}

	if in == nil {
		return out
	}

	out.Value = in.Value
	out.Currency = in.Currency

	return out
}
