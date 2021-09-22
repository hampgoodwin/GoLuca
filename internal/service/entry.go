package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/pagination"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func (s *Service) GetEntries(ctx context.Context, cursor string, limit string) ([]transaction.Entry, *string, error) {
	limitInt, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return nil, nil, errors.FlagWrap(
			err, errors.NotValidRequest,
			fmt.Sprintf("failed parsing provided limit query parameter %q", limit),
			"parsing limit query param")
	}
	limitInt++ // we always want one more than the size of the page, the extra at the end of the resultset serves as starting record for the next page
	var id string
	var createdAt time.Time
	if cursor != "" {
		id, createdAt, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, errors.FlagWrap(err, errors.NotValidRequest, err.Error(), "decoding base64 cursor")
		}
	}

	entries, err := s.repository.GetEntries(ctx, id, createdAt, limitInt)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting entries from database")
	}

	encodedCursor := ""
	if len(entries) == int(limitInt) {
		encodedCursor = pagination.EncodeCursor(entries[len(entries)-1].CreatedAt, entries[len(entries)-1].ID)
		entries = entries[:len(entries)-1]
	}
	return entries, &encodedCursor, nil
}
