package service

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrParseType = errors.New("error parsing string as type")

// Account represents a collection of entries into a logical grouping
type Account struct {
	ID        string    `validate:"required,uuid7"`
	ParentID  string    `validate:"omitempty,uuid7"`
	Name      string    `validate:"required"`
	Type      Type      `validate:"required"`
	Basis     Basis     `validate:"required"`
	CreatedAt time.Time `validate:"required"`
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
	// Slug should not be accessed by dependent code. It is exported for validation reasons
	Slug string `validate:"oneof=asset liability equity revenue expense gain loss"`
}

// safer enums for Type enum
var (
	TypeUnspecified = Type{""}
	TypeAsset       = Type{"asset"}
	TypeLiability   = Type{"liability"}
	TypeEquity      = Type{"equity"}
	TypeRevenue     = Type{"revenue"}
	TypeExpense     = Type{"expense"}
	TypeGain        = Type{"gain"}
	TypeLoss        = Type{"loss"}
)

// typeAsStringMap is used in parsing a string to a type
var typeAsStringMap = map[string]Type{
	"":          TypeUnspecified,
	"asset":     TypeAsset,
	"liability": TypeLiability,
	"equity":    TypeEquity,
	"revenue":   TypeRevenue,
	"expense":   TypeExpense,
	"gain":      TypeGain,
	"loss":      TypeLoss,
}

func ParseType(t string) Type {
	if v, ok := typeAsStringMap[t]; ok {
		return v
	}
	return TypeUnspecified
}

func (t Type) String() string {
	return t.Slug
}

type Basis struct {
	// Slug should not be accessed by dependent code. It is exported for validation reasons
	Slug string `validate:"oneof=debit credit"`
}

var (
	BasisUnspecified = Basis{""}
	BasisDebit       = Basis{"debit"}
	BasisCredit      = Basis{"credit"}
)

// typeAsStringMap is used in parsing a string to a type
var basisAsStringMap = map[string]Basis{
	"":       BasisUnspecified,
	"debit":  BasisDebit,
	"credit": BasisCredit,
}

func ParseBasis(b string) Basis {
	if v, ok := basisAsStringMap[b]; ok {
		return v
	}
	return BasisUnspecified
}

func (b Basis) String() string {
	return b.Slug
}
