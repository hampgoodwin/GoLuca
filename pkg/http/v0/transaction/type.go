package transaction

import (
	httpamount "github.com/hampgoodwin/GoLuca/pkg/http/v0/amount"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

type CreateTransaction struct {
	Description string        `json:"description" validate:"required"`
	Entries     []CreateEntry `json:"entries,omitempty" validate:"dive,gte=1"`
}

func (t CreateTransaction) IsZero() bool {
	if t.Description != "" {
		return false
	}
	if t.Entries != nil {
		return false
	}
	return true
}

type CreateEntry struct {
	Description   string            `json:"description"`
	DebitAccount  string            `json:"debitAccount" validate:"required,KSUID"`
	CreditAccount string            `json:"creditAccount" validate:"required,KSUID"`
	Amount        httpamount.Amount `json:"amount" validate:"required"`
}

type Transaction transaction.Transaction

type Entry transaction.Entry
