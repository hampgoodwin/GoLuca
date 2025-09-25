package account

import (
	accountv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/account/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAccountFromProtoCreateAccount(in *accountv1.CreateAccountRequest) Account {
	out := Account{}

	if in == nil {
		return out
	}

	out.ParentID = in.GetParentId()
	out.Name = in.Name
	out.Type = pbAccountTypeToAccountType(in.Type)
	out.Basis = pbAccountBasisToAccountBasis(in.Basis)

	return out
}

func NewProtoAccountFromAccount(in Account) *accountv1.Account {
	if in == (Account{}) {
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

var accountTypeToPBAccountTypeMap = map[Type]accountv1.AccountType{
	TypeUnspecified: accountv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED,
	TypeAsset:       accountv1.AccountType_ACCOUNT_TYPE_ASSET,
	TypeLiability:   accountv1.AccountType_ACCOUNT_TYPE_LIABILITY,
	TypeEquity:      accountv1.AccountType_ACCOUNT_TYPE_EQUITY,
	TypeRevenue:     accountv1.AccountType_ACCOUNT_TYPE_REVENUE,
	TypeExpense:     accountv1.AccountType_ACCOUNT_TYPE_EXPENSE,
	TypeGain:        accountv1.AccountType_ACCOUNT_TYPE_GAIN,
	TypeLoss:        accountv1.AccountType_ACCOUNT_TYPE_LOSS,
}

func accountTypeToPBAccountType(in Type) accountv1.AccountType {
	if v, ok := accountTypeToPBAccountTypeMap[in]; ok {
		return v
	}
	return accountv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED
}

var pbAccountTypeToAccountTypeMap = map[accountv1.AccountType]Type{
	accountv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED: TypeUnspecified,
	accountv1.AccountType_ACCOUNT_TYPE_ASSET:       TypeAsset,
	accountv1.AccountType_ACCOUNT_TYPE_LIABILITY:   TypeLiability,
	accountv1.AccountType_ACCOUNT_TYPE_EQUITY:      TypeEquity,
	accountv1.AccountType_ACCOUNT_TYPE_REVENUE:     TypeRevenue,
	accountv1.AccountType_ACCOUNT_TYPE_EXPENSE:     TypeExpense,
	accountv1.AccountType_ACCOUNT_TYPE_GAIN:        TypeGain,
	accountv1.AccountType_ACCOUNT_TYPE_LOSS:        TypeLoss,
}

func pbAccountTypeToAccountType(in accountv1.AccountType) Type {
	if v, ok := pbAccountTypeToAccountTypeMap[in]; ok {
		return v
	}
	return TypeUnspecified
}

var accountBasisToPBAccountBasisMap = map[Basis]accountv1.Basis{
	BasisUnspecified: accountv1.Basis_BASIS_UNSPECIFIED,
	BasisDebit:       accountv1.Basis_BASIS_DEBIT,
	BasisCredit:      accountv1.Basis_BASIS_CREDIT,
}

func accountBasisToPBAccountBasis(in Basis) accountv1.Basis {
	if v, ok := accountBasisToPBAccountBasisMap[in]; ok {
		return v
	}
	return accountv1.Basis_BASIS_UNSPECIFIED
}

var pbAccountBasisToAccountBasisMap = map[accountv1.Basis]Basis{
	accountv1.Basis_BASIS_UNSPECIFIED: BasisUnspecified,
	accountv1.Basis_BASIS_DEBIT:       BasisDebit,
	accountv1.Basis_BASIS_CREDIT:      BasisCredit,
}

func pbAccountBasisToAccountBasis(in accountv1.Basis) Basis {
	if v, ok := pbAccountBasisToAccountBasisMap[in]; ok {
		return v
	}
	return BasisUnspecified
}
