package transaction

import (
	"fmt"
)

// Transaction ...
type Transaction struct {
	Description string  `validate:"required"`
	Entries     []Entry `validate:"required,dive,gte=0"`
}

// Entry ...
type Entry struct {
	// Account account.Account `validate:"required"`
	ID      int64   `validate:"gte=1"`
	Account string  `validate:"required"`
	Amount  float64 `validate:"gt=0"`
}

func (t Transaction) String() string {
	stringer := fmt.Sprintf(`%s\n`, t.Description)
	for _, event := range t.Entries {
		stringer += fmt.Sprintf("%s", event)
	}
	return stringer
}

func (e Entry) String() string {
	return fmt.Sprintf(`%s
%f
`, e.Account, e.Amount)
}
