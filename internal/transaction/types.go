package transaction

import (
	"time"

	"github.com/hampgoodwin/GoLuca/pkg/amount"
)

// Transaction ...
type Transaction struct {
	ID          string    `validate:"required,KSUID"`
	Description string    `validate:"required"`
	Entries     []Entry   `validate:"dive,gte=1"`
	CreatedAt   time.Time `validate:"required"`
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
