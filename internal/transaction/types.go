package transaction

import (
	"time"
)

// Transaction ...
type Transaction struct {
	ID          string    `validate:"required,uuid7"`
	Description string    `validate:"required"`
	Entries     []Entry   `validate:"required,gt=0,dive"`
	CreatedAt   time.Time `validate:"required"`
}

type Amount struct {
	// Value denotes the number of currency in an amount. The last two digits
	// are the cents of the given currency. While this value is of type
	// the limitation is 9223372036854775807 due to database type limitations
	Value    int64  `validate:"gte=0"`
	Currency string `validate:"len=3,alpha"`
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
	DebitAccount  string    `validate:"required,uuid7"`
	CreditAccount string    `validate:"required,uuid7"`
	Amount        Amount    `validate:"required"`
	CreatedAt     time.Time `validate:"required"`
}
