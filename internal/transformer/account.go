package transformer

import (
	modelv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/model/v1"
	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/account"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	httpaccount "github.com/hampgoodwin/GoLuca/pkg/http/v0/account"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAccountFromHTTPCreateAccount(in httpaccount.CreateAccount) account.Account {
	out := account.Account{}

	if in == (httpaccount.CreateAccount{}) {
		return out
	}

	out.ParentID = in.ParentID
	out.Name = in.Name
	out.Type = account.ParseType(in.Type)
	out.Basis = account.ParseBasis(in.Basis)

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
	out.Type = in.Type.String()
	out.Basis = in.Basis.String()
	out.CreatedAt = in.CreatedAt

	return out
}

func NewAccountFromRepoAccount(in repository.Account) account.Account {
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

func NewRepoAccountFromAccount(in account.Account) repository.Account {
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

func NewProtoAccountFromAccount(in account.Account) *modelv1.Account {
	if in == (account.Account{}) {
		return nil
	}

	out := &modelv1.Account{
		Id:        in.ID,
		ParentId:  nil, // filled in later based on zero-value
		Name:      in.Name,
		Type:      accountTypeToPBAccountType(in.Type),
		Basis:     accountBasisToPBAccountBasis(in.Basis),
		CreatedAt: timestamppb.New(in.CreatedAt),
	}

	if in.ParentID != "" {
		out.ParentId = &in.ParentID
	}

	return out
}

func NewAccountFromProtoCreateAccount(in *servicev1.CreateAccountRequest) account.Account {
	out := account.Account{}

	if in == nil {
		return out
	}

	out.ParentID = in.GetParentId()
	out.Name = in.Name
	out.Type = pbAccountTypeToAccountType(in.Type)
	out.Basis = pbAccountBasisToAccountBasis(in.Basis)

	return out
}
