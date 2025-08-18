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
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateAccount(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	nc, _ := nats.Connect(nats.DefaultURL)
	service := service.NewService(repository, nc)
	controller := NewController(service)
	s.SetGRPC(t, controller)

	req := &servicev1.CreateAccountRequest{
		Name:  "cash",
		Type:  modelv1.AccountType_ACCOUNT_TYPE_ASSET,
		Basis: modelv1.Basis_BASIS_DEBIT,
	}

	res, err := s.GRPCTestClient.CreateAccount(s.Ctx, req)
	s.Is.NoErr(err)

	s.Is.True(res != nil)

	s.Is.True(res.GetAccount().GetId() != "")
	s.Is.True(res.GetAccount().GetParentId() == "")
	s.Is.Equal(res.GetAccount().GetName(), req.Name)
	s.Is.Equal(res.GetAccount().GetType(), req.Type)
	s.Is.Equal(res.GetAccount().GetBasis(), req.Basis)
	s.Is.True(res.GetAccount().GetCreatedAt() != nil)
}

func TestCreateAccount_InvalidRequestBody(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	nc, _ := nats.Connect(nats.DefaultURL)
	service := service.NewService(repository, nc)
	controller := NewController(service)
	s.SetGRPC(t, controller)

	req := &servicev1.CreateAccountRequest{}

	res, err := s.GRPCTestClient.CreateAccount(s.Ctx, req)
	s.Is.True(res == nil)
	s.Is.True(err != nil)

	st := status.Convert(err)
	s.Is.True(st.Code() == codes.InvalidArgument)

	fieldViolations := map[string]struct{}{
		"CreateAccountRequest.Name":  {},
		"CreateAccountRequest.Type":  {},
		"CreateAccountRequest.Basis": {},
	}
	for _, detail := range st.Details() {
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

func TestGetAccount(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	nc, _ := nats.Connect(nats.DefaultURL)
	service := service.NewService(repository, nc)
	controller := NewController(service)
	s.SetGRPC(t, controller)

	createAccountRequest := &servicev1.CreateAccountRequest{
		Name:  "cash",
		Type:  modelv1.AccountType_ACCOUNT_TYPE_ASSET,
		Basis: modelv1.Basis_BASIS_DEBIT,
	}

	createAccountResponse, err := s.GRPCTestClient.CreateAccount(s.Ctx, createAccountRequest)
	s.Is.NoErr(err)
	createdAccount := createAccountResponse.GetAccount()

	getAccountRequest := &servicev1.GetAccountRequest{
		AccountId: createdAccount.GetId(),
	}
	getAccountResponse, err := s.GRPCTestClient.GetAccount(s.Ctx, getAccountRequest)
	s.Is.NoErr(err)

	s.Is.Equal(getAccountResponse.GetAccount().GetId(), createdAccount.GetId())
	s.Is.Equal(getAccountResponse.GetAccount().GetParentId(), createdAccount.GetParentId())
	s.Is.Equal(getAccountResponse.GetAccount().GetName(), createdAccount.GetName())
	s.Is.Equal(getAccountResponse.GetAccount().GetType(), createdAccount.GetType())
	s.Is.Equal(getAccountResponse.GetAccount().GetBasis(), createdAccount.GetBasis())
	s.Is.True(getAccountResponse.GetAccount().GetCreatedAt() != nil)
}

func TestGetAccount_InvalidRequestBody(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	nc, _ := nats.Connect(nats.DefaultURL)
	service := service.NewService(repository, nc)
	controller := NewController(service)
	s.SetGRPC(t, controller)

	getAccountRequest := &servicev1.GetAccountRequest{
		AccountId: "not a valid uuidv7",
	}
	getAccountResponse, err := s.GRPCTestClient.GetAccount(s.Ctx, getAccountRequest)
	s.Is.True(getAccountResponse == nil)
	s.Is.True(err != nil)

	st := status.Convert(err)
	s.Is.True(st.Code() == codes.InvalidArgument)

	fieldViolations := map[string]struct{}{
		"GetAccountRequest.AccountId": {},
	}
	for _, detail := range st.Details() {
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

func TestListAccount(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	nc, _ := nats.Connect(nats.DefaultURL)
	service := service.NewService(repository, nc)
	controller := NewController(service)
	s.SetGRPC(t, controller)

	createAccountRequest := &servicev1.CreateAccountRequest{
		Name:  "cash",
		Type:  modelv1.AccountType_ACCOUNT_TYPE_ASSET,
		Basis: modelv1.Basis_BASIS_DEBIT,
	}
	_, err := s.GRPCTestClient.CreateAccount(s.Ctx, createAccountRequest)
	s.Is.NoErr(err)
	createAccountRequest.Name += " another"
	_, err = s.GRPCTestClient.CreateAccount(s.Ctx, createAccountRequest)
	s.Is.NoErr(err)
	createAccountRequest.Name += " another"
	_, err = s.GRPCTestClient.CreateAccount(s.Ctx, createAccountRequest)
	s.Is.NoErr(err)
	createAccountRequest.Name += " another"
	_, err = s.GRPCTestClient.CreateAccount(s.Ctx, createAccountRequest)
	s.Is.NoErr(err)

	listAccountsRequest := &servicev1.ListAccountsRequest{
		PageSize: 10,
	}
	listAccountsResponse, err := s.GRPCTestClient.ListAccounts(s.Ctx, listAccountsRequest)
	s.Is.NoErr(err)

	s.Is.True(len(listAccountsResponse.GetAccounts()) == 4)
}

func TestListAccount_InvalidRequestBody(t *testing.T) {
	s := test.GetScope(t)
	repository := repository.NewRepository(s.DB)
	nc, _ := nats.Connect(nats.DefaultURL)
	service := service.NewService(repository, nc)
	controller := NewController(service)
	s.SetGRPC(t, controller)

	listAccountsRequest := &servicev1.ListAccountsRequest{
		PageSize:  0,
		PageToken: "invalid token",
	}
	listAccountsResponse, err := s.GRPCTestClient.ListAccounts(s.Ctx, listAccountsRequest)
	s.Is.True(listAccountsResponse == nil)
	s.Is.True(err != nil)

	st := status.Convert(err)
	s.Is.True(st.Code() == codes.InvalidArgument)
	s.Is.True(st.Message() == fmt.Sprintf("invalid token %q", "invalid token"))
}
