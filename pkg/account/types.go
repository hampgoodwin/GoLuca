package account

import (
	"fmt"
)

// Account represents a collection of entries into a logical grouping
type Account struct {
	Parent *Account
	ID     uint32 `validate:"required,gte=0"`
	Name   string `validate:"required"`
	Type   Type   `validate:"required,gt=0,lte=6"`
	Basis  string `validate:"required,eq=debit,credit"`
}

func (a Account) String() string {
	stringer := fmt.Sprintf(`ID: %d
	Name: %s
	Type: %s
	Basis: %s
	Parent Account Name: %s`,
		a.ID, a.Name, a.Type, a.Basis, a.Parent.Name)
	return stringer
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
	Lose
)

func (t Type) String() string {
	return [...]string{"Asset", "Liability", "Equity", "Revenue", "Expense", "Gain", "Lose"}[t]
}
