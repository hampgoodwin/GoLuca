package transformer

import (
	"strconv"

	"github.com/hampgoodwin/GoLuca/internal/http/v0/api"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/amount"
	"github.com/hampgoodwin/errors"
)

func NewAmountFromHTTPAmount(a api.Amount) (amount.Amount, error) {
	out := amount.Amount{}
	if a == (api.Amount{}) {
		return out, nil
	}
	value, err := strconv.ParseInt(a.Value, 10, 64)
	if err != nil {
		return out, errors.Wrap(err, "converting amount string value to integer64")
	}
	out.Value = value
	out.Currency = a.Currency

	if err := validate.Validate(out); err != nil {
		return out, err
	}

	return out, nil
}
