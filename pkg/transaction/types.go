package transaction

import (
	"time"

	"github.com/hampgoodwin/GoLuca/pkg/amount"
)

// Transaction ...
type Transaction struct {
	ID          string    `json:"id" validate:"required,uuid4"`
	Description string    `json:"description" validate:"required"`
	Entries     []Entry   `json:"entries,omitempty" validate:"dive,gte=1"`
	CreatedAt   time.Time `json:"createdAt" validate:"required"`
}

// Entry ...
type Entry struct {
	ID            string        `json:"id" validate:"required,uuid4"`
	TransactionID string        `json:"transaction_id" validate:"required,uuid4"`
	Description   string        `json:"description"`
	DebitAccount  string        `json:"debitAccount" validate:"required,uuid4"`
	CreditAccount string        `json:"creditAccount" validate:"required,uuid4"`
	Amount        amount.Amount `json:"amount" validate:"required"`
	CreatedAt     time.Time     `json:"createdAt" validate:"required"`
}
