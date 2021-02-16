package account

import (
	"encoding/json"
)

// Account represents a collection of entries into a logical grouping
type Account struct {
	ID       int64  `json:"id" validate:"gte=0"`
	ParentID int64  `json:"parent_id" validate:"gte=0"`
	Name     string `json:"name" validate:"required"`
	Type     Type   `json:"type" validate:"required,gt=0,lte=7"`
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
type Type int

// iota const's for account Type
const (
	Asset = iota
	Liability
	Equity
	Revenue
	Expense
	Gain
	Loss
)

func (t Type) String() string {
	return [...]string{"Asset", "Liability", "Equity", "Revenue", "Expense", "Gain", "Loss"}[t-1]
}
