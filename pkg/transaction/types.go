package transaction

import (
	"encoding/json"
	"fmt"
)

// Transaction ...
type Transaction struct {
	ID          string  `json:"id" validate:"uuid4"`
	Description string  `json:"description" validate:"required"`
	Entries     []Entry `json:"entries,omitempty" validate:"dive,gte=2"`
}

// Entry ...
type Entry struct {
	// Account account.Account `validate:"required"`
	ID            string  `json:"id" validate:"uuid4"`
	TransactionID string  `json:"transaction_id" validate:"uudi4"`
	AccountID     int64   `json:"account_id"  validate:"required,gte=0"`
	Amount        float64 `json:"amount" validate:"required,ne=0"`
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
