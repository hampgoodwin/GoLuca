package controller

import (
	"context"
	"fmt"

	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/errors"
)

func (c *Controller) GetAccount(ctx context.Context, req *servicev1.GetAccountRequest) (*servicev1.GetAccountResponse, error) {
	c.log.Info(fmt.Sprintf("getting account %q", req.AccountId)) // TODO factor out in a gRPC router logger for incoming/outgoing requests.

	serviceAccount, err := c.service.GetAccount(ctx, req.AccountId)
	if err != nil { // TODO update to use status package
		return nil, errors.WithMessage(err, "getting account")
	}

	account := transformer.NewPBAccountFromAccount(serviceAccount)
	// I believe there is a way to validate struct by struct, so we'll have to use helper
	// functions to set validations on a per-field basis for protobuf structs, and call
	// those helpers
	// TODO: Validate

	return &servicev1.GetAccountResponse{Account: account}, nil
}

func (c *Controller) ListAccounts(ctx context.Context, req *servicev1.ListAccountsRequest) (*servicev1.ListAccountsResponse, error) {

	return &servicev1.ListAccountsResponse{}, nil
}

func (c *Controller) CreateAccount(ctx context.Context, create *servicev1.CreateAccountRequest) (*servicev1.CreateAccountResponse, error) {
	c.log.Info("creating account") // TODO factor out in a gRPC router logger for incoming/outgoing requests.

	serviceAccount := transformer.NewAccountFromPBCreateAccount(create)

	serviceAccount, err := c.service.CreateAccount(ctx, serviceAccount)
	if err != nil {
		return nil, errors.WithMessage(err, "creating account")
	}

	account := transformer.NewPBAccountFromAccount(serviceAccount)
	return &servicev1.CreateAccountResponse{Account: account}, nil
}
