package transformer

import (
	"github.com/hampgoodwin/GoLuca/pkg/account"
	httpaccount "github.com/hampgoodwin/GoLuca/pkg/http/v0/account"
)

func NewAccountFromHTTPCreateAccount(in httpaccount.CreateAccount) account.Account {
	out := account.Account{}

	if in == (httpaccount.CreateAccount{}) {
		return out
	}

	out.ParentID = in.ParentID
	out.Name = in.Name
	out.Type = in.Type
	out.Basis = in.Basis

	return out
}

func NewHTTPAccountFromAccount(in account.Account) httpaccount.Account {
	out := httpaccount.Account{}

	if in == (account.Account{}) {
		return out
	}

	out.ID = in.ID
	out.ParentID = in.ParentID
	out.Name = in.Name
	out.Type = in.Type
	out.Basis = in.Basis
	out.CreatedAt = in.CreatedAt

	return out
}
