package repository

import (
	"time"
)

type Account struct {
	ID        string
	ParentID  string
	Name      string
	Type      string
	Basis     string
	CreatedAt time.Time
}

type Transaction struct {
	ID          string
	Description string
	Entries     []Entry
	CreatedAt   time.Time
}

func (t Transaction) IsZero() bool {
	if t.Description != "" {
		return false
	}
	if t.Entries != nil {
		return false
	}
	return true
}

type Entry struct {
	ID            string
	TransactionID string
	Description   string
	DebitAccount  string
	CreditAccount string
	Amount        Amount
	CreatedAt     time.Time
}

type Amount struct {
	Value    int64  `validate:"gte=0"`
	Currency string `validate:"len=3,alpha"`
}
