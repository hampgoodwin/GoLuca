package account

import (
	"encoding/json"
)

// Account represents a collection of entries into a logical grouping
type Account struct {
	ID       int64  `json:"id" validate:"uuid4"`
	ParentID int64  `json:"parentId" validate:"uuid4"`
	Name     string `json:"name" validate:"required"`
	Type     Type   `json:"type" validate:"required,oneof=asset liability equity revenue expense gain loss"`
	Basis    string `json:"basis" validate:"required,oneof=debit credit"`
}

func (a Account) String() string {
	b, err := json.Marshal(a)
	if err != nil {
		return "err"
	}
	return string(b)
}

// Type represents the type of account
type Type string

// iota const's for account Type
const (
	Asset     Type = "asset"
	Liability Type = "liablity"
	Equity    Type = "equity"
	Revenue   Type = "revenue"
	Expense   Type = "expense"
	Gain      Type = "gain"
	Loss      Type = "loss"
)

func (t Type) String() string {
	return string(t)
}
