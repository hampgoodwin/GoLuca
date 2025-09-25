package service

import (
	accountv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/account/v1"
	"github.com/hampgoodwin/GoLuca/internal/account/repository"
)

func newAccountFromRepoAccount(in repository.Account) Account {
	out := Account{}

	if in == (repository.Account{}) {
		return out
	}

	out.ID = in.ID
	out.ParentID = in.ParentID
	out.Name = in.Name
	out.Type = ParseType(in.Type)
	out.Basis = ParseBasis(in.Basis)
	out.CreatedAt = in.CreatedAt

	return out
}

func NewRepoAccountFromAccount(in Account) repository.Account {
	out := repository.Account{}

	if in == (Account{}) {
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
