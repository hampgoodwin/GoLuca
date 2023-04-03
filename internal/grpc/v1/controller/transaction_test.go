package controller

import (
	"fmt"
	"testing"

	modelv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/model/v1"
	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/hampgoodwin/GoLuca/internal/test"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateTransaction(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	nc, _ := nats.Connect(nats.DefaultURL)
	nec, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	service := service.NewService(s.Env.Log, repository, nec)
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

func TestCreateTransaction_InvalidRequestBody(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	nc, _ := nats.Connect(nats.DefaultURL)
	nec, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	service := service.NewService(s.Env.Log, repository, nec)
	controller := NewController(s.Env.Log, service)
	s.SetGRPC(t, controller)

	createTransactionRequest := &servicev1.CreateTransactionRequest{
		Entries: []*servicev1.CreateEntry{{Amount: &modelv1.Amount{}}},
	}
	createTransactionResponse, err := s.GRPCTestClient.CreateTransaction(s.Ctx, createTransactionRequest)
	s.Is.True(err != nil)
	s.Is.True(createTransactionResponse == nil)

	st := status.Convert(err)

	s.Is.True(st.Code() == codes.InvalidArgument)

	fieldViolations := map[string]struct{}{
		"CreateTransactionRequest.Description":                {},
		"CreateTransactionRequest.Entries[0].Description":     {},
		"CreateTransactionRequest.Entries[0].DebitAccount":    {},
		"CreateTransactionRequest.Entries[0].CreditAccount":   {},
		"CreateTransactionRequest.Entries[0].Amount.Currency": {},
	}
	details := st.Details()
	for _, detail := range details {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range t.GetFieldViolations() {
				if _, ok := fieldViolations[violation.Field]; !ok {
					fmt.Printf("%s not found in expected field violations", violation.Field)
					s.Is.Fail()
				}
				delete(fieldViolations, violation.Field)
			}
		}
	}
	if len(fieldViolations) != 0 {
		fmt.Println(fieldViolations)
	}
	s.Is.True(len(fieldViolations) == 0)
}

func TestGetTransaction(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	nc, _ := nats.Connect(nats.DefaultURL)
	nec, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	service := service.NewService(s.Env.Log, repository, nec)
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

	// NOW GET IT
	getTransactionRequest := &servicev1.GetTransactionRequest{
		TransactionId: createTransactionResponse.GetTransaction().GetId(),
	}
	getTransactionResponse, err := s.GRPCTestClient.GetTransaction(s.Ctx, getTransactionRequest)
	s.Is.NoErr(err)

	s.Is.True(getTransactionResponse != nil)
	createdTransaction := createTransactionResponse.GetTransaction()
	transaction := getTransactionResponse.GetTransaction()
	s.Is.Equal(transaction, createdTransaction)
}

func TestListTransactions(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	nc, _ := nats.Connect(nats.DefaultURL)
	nec, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	service := service.NewService(s.Env.Log, repository, nec)
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

	// NOW GET IT
	listTransactionRequest := &servicev1.ListTransactionsRequest{
		PageSize:  0,
		PageToken: "",
	}
	listTransactionsResponse, err := s.GRPCTestClient.ListTransactions(s.Ctx, listTransactionRequest)
	s.Is.NoErr(err)

	s.Is.True(listTransactionsResponse != nil)
	createdTransaction := createTransactionResponse.GetTransaction()
	transaction := listTransactionsResponse.GetTransactions()[0]
	s.Is.Equal(transaction, createdTransaction)
}
