package account

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrParseType = errors.New("error parsing string as type")

// Account represents a collection of entries into a logical grouping
type Account struct {
	ID        string    `json:"id" validate:"required,uuid4"`
	ParentID  string    `json:"parentId,omitempty" validate:"omitempty,uuid4"`
	Name      string    `json:"name" validate:"required"`
	Type      Type      `json:"type" validate:"required"`
	Basis     string    `json:"basis" validate:"required,oneof=debit credit"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}

func (a Account) String() string {
	b, err := json.Marshal(a)
	if err != nil {
		return "err"
	}
	return string(b)
}

// Type represents the type of account
type Type struct {
	slug string `validate:"required,oneof=asset liability equity revenue expense gain loss"`
}

// iota const's for account Type
var (
	TypeUnknown   = Type{""}
	TypeAsset     = Type{"asset"}
	TypeLiability = Type{"liability"}
	TypeEquity    = Type{"equity"}
	TypeRevenue   = Type{"revenue"}
	TypeExpense   = Type{"expense"}
	TypeGain      = Type{"gain"}
	TypeLoss      = Type{"loss"}
)

// typeAsStringMap is used in parsing a string to a type
var typeAsStringMap = map[string]Type{
	"":         TypeUnknown,
	"asset":    TypeAsset,
	"liablity": TypeLiability,
	"equity":   TypeEquity,
	"revenue":  TypeRevenue,
	"expense":  TypeExpense,
	"gain":     TypeGain,
	"loss":     TypeLoss,
}

func ParseType(t string) Type {
	if v, ok := typeAsStringMap[t]; ok {
		return v
	}
	return TypeUnknown
}

func (t Type) String() string {
	return t.slug
}
