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
