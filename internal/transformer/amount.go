package transformer

import (
	"strconv"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/httpapi"
	"github.com/hampgoodwin/GoLuca/pkg/amount"
)

func NewAmountFromHTTPAmount(a httpapi.Amount) (amount.Amount, error) {
	out := amount.Amount{}
	value, err := strconv.ParseInt(a.Value, 10, 64)
	if err != nil {
		return out, errors.Wrap(err, "converting amount string value to integer64")
	}
	out.Value = value
	out.Currency = a.Currency
	return out, nil
}
