package account

import "github.com/hampgoodwin/GoLuca/pkg/account"

type CreateAccount struct {
	ParentID string       `json:"parentId,omitempty" validate:"omitempty,uuid4"`
	Name     string       `json:"name" validate:"required"`
	Type     account.Type `json:"type" validate:"required,oneof=asset liability equity revenue expense gain loss"`
	Basis    string       `json:"basis" validate:"required,oneof=debit credit"`
}

type Account account.Account
