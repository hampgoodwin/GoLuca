package service

import (
	"github.com/hampgoodwin/GoLuca/internal/account"
	"github.com/hampgoodwin/GoLuca/internal/account/repository"
)

func newAccountFromRepoAccount(in repository.Account) account.Account {
	out := account.Account{}

	if in == (repository.Account{}) {
		return out
	}

	out.ID = in.ID
	out.ParentID = in.ParentID
	out.Name = in.Name
	out.Type = account.ParseType(in.Type)
	out.Basis = account.ParseBasis(in.Basis)
	out.CreatedAt = in.CreatedAt

	return out
}

func newRepoAccountFromAccount(in account.Account) repository.Account {
	out := repository.Account{}

	if in == (account.Account{}) {
		return out
	}

	out.ID = in.ID
	out.ParentID = in.ParentID
	out.Name = in.Name
	out.Type = in.Type.String()
	out.Basis = in.Basis.String()
	out.CreatedAt = in.CreatedAt

	return out
}
