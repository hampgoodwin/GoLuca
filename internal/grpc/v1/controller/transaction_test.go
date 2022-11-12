package controller

import (
	"testing"

	modelv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/model/v1"
	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/hampgoodwin/GoLuca/internal/test"
)

func TestCreateTransaction(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	service := service.NewService(s.Env.Log, repository)
	controller := NewController(s.Env.Log, service)
	s.SetGRPC(t, controller)

	createDebitAccountRequest := &servicev1.CreateAccountRequest{
		Name:  "debit",
		Type:  modelv1.AccountType_ACCOUNT_TYPE_ASSET,
		Basis: modelv1.Basis_BASIS_DEBIT,
	}
	assetAccount, err := s.GRPCTestClient.CreateAccount(s.Ctx, createDebitAccountRequest)
	s.Is.NoErr(err)
	createEquityAccountRequest := &servicev1.CreateAccountRequest{
		Name:  "liability",
		Type:  modelv1.AccountType_ACCOUNT_TYPE_EQUITY,
		Basis: modelv1.Basis_BASIS_CREDIT,
	}
	equityAccount, err := s.GRPCTestClient.CreateAccount(s.Ctx, createEquityAccountRequest)
	s.Is.NoErr(err)

	createTransactionRequest := &servicev1.CreateTransactionRequest{
		Description: "test",
		Entries: []*servicev1.CreateEntry{{
			Description:   "entry1",
			DebitAccount:  assetAccount.Account.Id,
			CreditAccount: equityAccount.Account.Id,
			Amount: &modelv1.Amount{
				Value:    2000,
				Currency: "USD",
			},
		}},
	}
	createTransactionResponse, err := s.GRPCTestClient.CreateTransaction(s.Ctx, createTransactionRequest)
	s.Is.NoErr(err)
	s.Is.True(createTransactionResponse != nil)

	s.Is.True(len(createTransactionResponse.Transaction.Entries) == 1)
	entry := createTransactionResponse.Transaction.Entries[0]
	s.Is.True(entry.DebitAccount == assetAccount.Account.Id)
	s.Is.True(entry.CreditAccount == equityAccount.Account.Id)
	s.Is.True(entry.Description == createTransactionRequest.Entries[0].Description)
	s.Is.True(entry.TransactionId == createTransactionResponse.Transaction.Id)
	amount := entry.Amount
	s.Is.True(amount.Value == createTransactionRequest.Entries[0].Amount.Value)
	s.Is.True(amount.Currency == createTransactionRequest.Entries[0].Amount.Currency)
}
