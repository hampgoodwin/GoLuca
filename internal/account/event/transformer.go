package event

import (
	accountv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/account/v1"
	"github.com/hampgoodwin/GoLuca/internal/account"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAccountFromProtoCreateAccount(in *accountv1.CreateAccountRequest) account.Account {
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

func NewProtoAccountFromAccount(in account.Account) *accountv1.Account {
	if in == (account.Account{}) {
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

var accountTypeToPBAccountTypeMap = map[account.Type]accountv1.AccountType{
	account.TypeUnspecified: accountv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED,
	account.TypeAsset:       accountv1.AccountType_ACCOUNT_TYPE_ASSET,
	account.TypeLiability:   accountv1.AccountType_ACCOUNT_TYPE_LIABILITY,
	account.TypeEquity:      accountv1.AccountType_ACCOUNT_TYPE_EQUITY,
	account.TypeRevenue:     accountv1.AccountType_ACCOUNT_TYPE_REVENUE,
	account.TypeExpense:     accountv1.AccountType_ACCOUNT_TYPE_EXPENSE,
	account.TypeGain:        accountv1.AccountType_ACCOUNT_TYPE_GAIN,
	account.TypeLoss:        accountv1.AccountType_ACCOUNT_TYPE_LOSS,
}

func accountTypeToPBAccountType(in account.Type) accountv1.AccountType {
	if v, ok := accountTypeToPBAccountTypeMap[in]; ok {
		return v
	}
	return accountv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED
}

var pbAccountTypeToAccountTypeMap = map[accountv1.AccountType]account.Type{
	accountv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED: account.TypeUnspecified,
	accountv1.AccountType_ACCOUNT_TYPE_ASSET:       account.TypeAsset,
	accountv1.AccountType_ACCOUNT_TYPE_LIABILITY:   account.TypeLiability,
	accountv1.AccountType_ACCOUNT_TYPE_EQUITY:      account.TypeEquity,
	accountv1.AccountType_ACCOUNT_TYPE_REVENUE:     account.TypeRevenue,
	accountv1.AccountType_ACCOUNT_TYPE_EXPENSE:     account.TypeExpense,
	accountv1.AccountType_ACCOUNT_TYPE_GAIN:        account.TypeGain,
	accountv1.AccountType_ACCOUNT_TYPE_LOSS:        account.TypeLoss,
}

func pbAccountTypeToAccountType(in accountv1.AccountType) account.Type {
	if v, ok := pbAccountTypeToAccountTypeMap[in]; ok {
		return v
	}
	return account.TypeUnspecified
}

var accountBasisToPBAccountBasisMap = map[account.Basis]accountv1.Basis{
	account.BasisUnspecified: accountv1.Basis_BASIS_UNSPECIFIED,
	account.BasisDebit:       accountv1.Basis_BASIS_DEBIT,
	account.BasisCredit:      accountv1.Basis_BASIS_CREDIT,
}

func accountBasisToPBAccountBasis(in account.Basis) accountv1.Basis {
	if v, ok := accountBasisToPBAccountBasisMap[in]; ok {
		return v
	}
	return accountv1.Basis_BASIS_UNSPECIFIED
}

var pbAccountBasisToAccountBasisMap = map[accountv1.Basis]account.Basis{
	accountv1.Basis_BASIS_UNSPECIFIED: account.BasisUnspecified,
	accountv1.Basis_BASIS_DEBIT:       account.BasisDebit,
	accountv1.Basis_BASIS_CREDIT:      account.BasisCredit,
}

func pbAccountBasisToAccountBasis(in accountv1.Basis) account.Basis {
	if v, ok := pbAccountBasisToAccountBasisMap[in]; ok {
		return v
	}
	return account.BasisUnspecified
}
