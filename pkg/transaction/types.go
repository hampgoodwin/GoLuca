package transaction

import (
	"encoding/json"
	"fmt"
	"time"
)

// Transaction ...
type Transaction struct {
	ID          string    `json:"id" validate:"required,uuid4"`
	Description string    `json:"description" validate:"required"`
	Entries     []Entry   `json:"entries,omitempty" validate:"dive,gte=2"`
	CreatedAt   time.Time `json:"createdAt" validate:"required"`
}

// Entry ...
type Entry struct {
	// Account account.Account `validate:"required"`
	ID            string    `json:"id" validate:"required,uuid4"`
	TransactionID string    `json:"transaction_id" validate:"required,uuid4"`
	AccountID     string    `json:"accountId"  validate:"required,uuid4"`
	Amount        float64   `json:"amount" validate:"required,ne=0"`
	CreatedAt     time.Time `json:"createdAt" validate:"required"`
}

func (t Transaction) String() string {
	stringer := fmt.Sprintf(`%s\n`, t.Description)
	for _, event := range t.Entries {
		stringer += event.String()
	}
	return stringer
}

func (e Entry) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "err"
	}
	return string(b)
}
