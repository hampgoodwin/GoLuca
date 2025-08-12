package transaction

import (
	"time"

	"github.com/hampgoodwin/GoLuca/internal/amount"
)

// Transaction ...
type Transaction struct {
	ID          string    `validate:"required,uuid7"`
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
	ID            string `validate:"required,uuid7"`
	TransactionID string `validate:"required,uuid7"`
	Description   string
	DebitAccount  string        `validate:"required,uuid7"`
	CreditAccount string        `validate:"required,uuid7"`
	Amount        amount.Amount `validate:"required"`
	CreatedAt     time.Time     `validate:"required"`
}
