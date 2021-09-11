package service

import (
	"context"

	"github.com/hampgoodwin/GoLuca/internal/data"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func GetEntries(ctx context.Context, cursor int64, limit int64) ([]transaction.Entry, error) {
	entries, err := data.GetEntries(ctx, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting entries from database")
	}
	return entries, nil
}
