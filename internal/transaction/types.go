package transaction

import (
	"time"

	"github.com/hampgoodwin/GoLuca/internal/amount"
)

// Transaction ...
type Transaction struct {
	ID          string    `validate:"required,KSUID"`
	Description string    `validate:"required"`
	Entries     []Entry   `validate:"required,gt=0,dive"`
	CreatedAt   time.Time `validate:"required"`
}

func (t Transaction) IsZero() bool {
	if t.Description != "" {
		return false
	}
	if t.Entries != nil {
		return false
	}
	if t.ID != "" {
		return false
	}
	if t.CreatedAt != (time.Time{}) {
		return false
	}
	return true
}

// Entry ...
type Entry struct {
	ID            string `validate:"required,KSUID"`
	TransactionID string `validate:"required,KSUID"`
	Description   string
	DebitAccount  string        `validate:"required,KSUID"`
	CreditAccount string        `validate:"required,KSUID"`
	Amount        amount.Amount `validate:"required"`
	CreatedAt     time.Time     `validate:"required"`
}
