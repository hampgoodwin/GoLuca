package controller

import (
	"context"
	"fmt"

	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/errors"
)

func (c *Controller) GetAccount(ctx context.Context, accountID string) (*servicev1.GetAccountResponse, error) {
	c.log.Info(fmt.Sprintf("getting accout %q", accountID)) // TODO factor out in a gRPC router logger for incoming/outgoing requests.

	serviceAccount, err := c.service.GetAccount(ctx, accountID)
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
