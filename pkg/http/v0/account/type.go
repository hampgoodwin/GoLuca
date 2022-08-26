package account

import "time"

type CreateAccount struct {
	ParentID string `json:"parentId,omitempty" validate:"omitempty,uuid4"`
	Name     string `json:"name" validate:"required"`
	Type     string `json:"type" validate:"required,oneof=asset liability equity revenue expense gain loss"`
	Basis    string `json:"basis" validate:"required,oneof=debit credit"`
}

type Account struct {
	ID        string    `json:"id" validate:"required,uuid4"`
	ParentID  string    `json:"parentId,omitempty" validate:"omitempty,uuid4"`
	Name      string    `json:"name" validate:"required"`
	Type      string    `json:"type" validate:"required"`
	Basis     string    `json:"basis" validate:"required,oneof=debit credit"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}
