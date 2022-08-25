package transformer

import (
	"github.com/hampgoodwin/GoLuca/pkg/account"
	httpaccount "github.com/hampgoodwin/GoLuca/pkg/http/account"
)

func NewAccountFromHTTPAccount(in httpaccount.CreateAccount) (account.Account, error) {
	out := account.Account{}

	if in == (httpaccount.CreateAccount{}) {
		return out, nil
	}

	out.ParentID = in.ParentID
	out.Name = in.Name
	out.Type = in.Type
	out.Basis = in.Basis

	return out, nil
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
