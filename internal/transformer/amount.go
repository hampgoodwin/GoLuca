package transformer

import (
	"strconv"

	"github.com/hampgoodwin/GoLuca/internal/amount"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	httpamount "github.com/hampgoodwin/GoLuca/pkg/http/v0/amount"
	"github.com/hampgoodwin/errors"
)

func NewAmountFromHTTPAmount(in httpamount.Amount) (amount.Amount, error) {
	out := amount.Amount{}
	if in == (httpamount.Amount{}) {
		return out, nil
	}
	value, err := strconv.ParseInt(in.Value, 10, 64)
	if err != nil {
		return out, errors.Wrap(err, "converting amount string value to integer64")
	}
	out.Value = value
	out.Currency = in.Currency

	if err := validate.Validate(out); err != nil {
		return out, err
	}

	return out, nil
}

func NewHTTPAmountFromAmount(in amount.Amount) httpamount.Amount {
	out := httpamount.Amount{}

	if in == (amount.Amount{}) {
		return out
	}

	out.Value = strconv.FormatInt(in.Value, 10)
	out.Currency = in.Currency

	return out
}

func NewAmountFromRepoAmount(in repository.Amount) amount.Amount {
	out := amount.Amount{}

	if in == (repository.Amount{}) {
		return out
	}

	out.Value = in.Value
	out.Currency = in.Currency

	return out
}

func NewRepoAmountFromAmount(in amount.Amount) repository.Amount {
	out := repository.Amount{}

	if in == (amount.Amount{}) {
		return out
	}

	out.Value = in.Value
	out.Currency = in.Currency

	return out
}
