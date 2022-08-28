package transformer

import (
	modelv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/model/v1"
	"github.com/hampgoodwin/GoLuca/internal/account"
)

var accountTypeToPBAccountTypeMap = map[account.Type]modelv1.AccountType{
	account.TypeUnspecified: modelv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED,
	account.TypeAsset:       modelv1.AccountType_ACCOUNT_TYPE_ASSET,
	account.TypeLiability:   modelv1.AccountType_ACCOUNT_TYPE_LIABILITY,
	account.TypeEquity:      modelv1.AccountType_ACCOUNT_TYPE_EQUITY,
	account.TypeRevenue:     modelv1.AccountType_ACCOUNT_TYPE_REVENUE,
	account.TypeExpense:     modelv1.AccountType_ACCOUNT_TYPE_EXPENSE,
	account.TypeGain:        modelv1.AccountType_ACCOUNT_TYPE_GAIN,
	account.TypeLoss:        modelv1.AccountType_ACCOUNT_TYPE_LOSS,
}

func accountTypeToPBAccountType(in account.Type) modelv1.AccountType {
	if v, ok := accountTypeToPBAccountTypeMap[in]; ok {
		return v
	}
	return modelv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED
}

var accountBasisToPBAccountBasisMap = map[account.Basis]modelv1.Basis{
	account.BasisUnspecified: modelv1.Basis_BASIS_UNSPECIFIED,
	account.BasisDebit:       modelv1.Basis_BASIS_DEBIT,
	account.BasisCredit:      modelv1.Basis_BASIS_CREDIT,
}

func accountBasisToPBAccountBasis(in account.Basis) modelv1.Basis {
	if v, ok := accountBasisToPBAccountBasisMap[in]; ok {
		return v
	}
	return modelv1.Basis_BASIS_UNSPECIFIED
}
