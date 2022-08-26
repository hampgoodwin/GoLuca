package transaction

import (
	"time"

	"github.com/hampgoodwin/GoLuca/pkg/amount"
)

// Transaction ...
type Transaction struct {
	ID          string    `validate:"required,uuid4"`
	Description string    `validate:"required"`
	Entries     []Entry   `validate:"dive,gte=1"`
	CreatedAt   time.Time `validate:"required"`
}

// Entry ...
type Entry struct {
	ID            string `validate:"required,uuid4"`
	TransactionID string `validate:"required,uuid4"`
	Description   string
	DebitAccount  string        `validate:"required,uuid4"`
	CreditAccount string        `validate:"required,uuid4"`
	Amount        amount.Amount `validate:"required"`
	CreatedAt     time.Time     `validate:"required"`
}
