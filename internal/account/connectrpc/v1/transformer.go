package connect

import (
	accountv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/account/v1"
	"github.com/hampgoodwin/GoLuca/internal/account/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAccountFromProtoCreateAccount(in *accountv1.CreateAccountRequest) service.Account {
	out := service.Account{}

	if in == nil {
		return out
	}

	out.ParentID = in.GetParentId()
	out.Name = in.Name
	out.Type = pbAccountTypeToAccountType(in.Type)
	out.Basis = pbAccountBasisToAccountBasis(in.Basis)

	return out
}

func NewProtoAccountFromAccount(in service.Account) *accountv1.Account {
	if in == (service.Account{}) {
		return nil
	}

	out := &accountv1.Account{
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

var accountTypeToPBAccountTypeMap = map[service.Type]accountv1.AccountType{
	service.TypeUnspecified: accountv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED,
	service.TypeAsset:       accountv1.AccountType_ACCOUNT_TYPE_ASSET,
	service.TypeLiability:   accountv1.AccountType_ACCOUNT_TYPE_LIABILITY,
	service.TypeEquity:      accountv1.AccountType_ACCOUNT_TYPE_EQUITY,
	service.TypeRevenue:     accountv1.AccountType_ACCOUNT_TYPE_REVENUE,
	service.TypeExpense:     accountv1.AccountType_ACCOUNT_TYPE_EXPENSE,
	service.TypeGain:        accountv1.AccountType_ACCOUNT_TYPE_GAIN,
	service.TypeLoss:        accountv1.AccountType_ACCOUNT_TYPE_LOSS,
}

func accountTypeToPBAccountType(in service.Type) accountv1.AccountType {
	if v, ok := accountTypeToPBAccountTypeMap[in]; ok {
		return v
	}
	return accountv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED
}

var pbAccountTypeToAccountTypeMap = map[accountv1.AccountType]service.Type{
	accountv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED: service.TypeUnspecified,
	accountv1.AccountType_ACCOUNT_TYPE_ASSET:       service.TypeAsset,
	accountv1.AccountType_ACCOUNT_TYPE_LIABILITY:   service.TypeLiability,
	accountv1.AccountType_ACCOUNT_TYPE_EQUITY:      service.TypeEquity,
	accountv1.AccountType_ACCOUNT_TYPE_REVENUE:     service.TypeRevenue,
	accountv1.AccountType_ACCOUNT_TYPE_EXPENSE:     service.TypeExpense,
	accountv1.AccountType_ACCOUNT_TYPE_GAIN:        service.TypeGain,
	accountv1.AccountType_ACCOUNT_TYPE_LOSS:        service.TypeLoss,
}

func pbAccountTypeToAccountType(in accountv1.AccountType) service.Type {
	if v, ok := pbAccountTypeToAccountTypeMap[in]; ok {
		return v
	}
	return service.TypeUnspecified
}

var accountBasisToPBAccountBasisMap = map[service.Basis]accountv1.Basis{
	service.BasisUnspecified: accountv1.Basis_BASIS_UNSPECIFIED,
	service.BasisDebit:       accountv1.Basis_BASIS_DEBIT,
	service.BasisCredit:      accountv1.Basis_BASIS_CREDIT,
}

func accountBasisToPBAccountBasis(in service.Basis) accountv1.Basis {
	if v, ok := accountBasisToPBAccountBasisMap[in]; ok {
		return v
	}
	return accountv1.Basis_BASIS_UNSPECIFIED
}

var pbAccountBasisToAccountBasisMap = map[accountv1.Basis]service.Basis{
	accountv1.Basis_BASIS_UNSPECIFIED: service.BasisUnspecified,
	accountv1.Basis_BASIS_DEBIT:       service.BasisDebit,
	accountv1.Basis_BASIS_CREDIT:      service.BasisCredit,
}

func pbAccountBasisToAccountBasis(in accountv1.Basis) service.Basis {
	if v, ok := pbAccountBasisToAccountBasisMap[in]; ok {
		return v
	}
	return service.BasisUnspecified
}
