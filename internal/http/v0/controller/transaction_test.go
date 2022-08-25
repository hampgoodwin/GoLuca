package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/hampgoodwin/GoLuca/internal/http/v0/api"
	"github.com/hampgoodwin/GoLuca/internal/test"
	"github.com/hampgoodwin/GoLuca/pkg/account"
	"github.com/hampgoodwin/GoLuca/pkg/amount"
	httpaccount "github.com/hampgoodwin/GoLuca/pkg/http/account"
)

func TestCreateTransaction(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	aReq := accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "cash",
			Type:  account.Asset,
			Basis: "debit",
		},
	}
	res := createAccount(t, &s, aReq)
	defer res.Body.Close()
	var aRes accountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.Is.NoErr(err)
	cashAccount := aRes.Account

	aReq = accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "revenue",
			Type:  account.Asset,
			Basis: "credit",
		},
	}
	res2 := createAccount(t, &s, aReq)
	defer res2.Body.Close()
	err = json.NewDecoder(res2.Body).Decode(&aRes)
	s.Is.NoErr(err)
	revenueAccount := aRes.Account

	tReq := transactionRequest{
		Transaction: api.Transaction{
			Description: "test",
			Entries: []api.Entry{
				{
					Description:   "",
					DebitAccount:  cashAccount.ID,
					CreditAccount: revenueAccount.ID,
					Amount: api.Amount{
						Value:    "100",
						Currency: "USD",
					},
				},
			},
		},
	}

	res3 := createTransaction(t, &s, tReq)
	defer res3.Body.Close()

	var tRes transactionResponse
	err = json.NewDecoder(res3.Body).Decode(&tRes)
	s.Is.NoErr(err)

	s.Is.True(tRes != (transactionResponse{}))

	s.Is.True(tRes.Entries[0].CreditAccount == revenueAccount.ID)
	s.Is.True(tRes.Entries[0].DebitAccount == cashAccount.ID)
	s.Is.True(tRes.Entries[0].Amount == amount.Amount{Value: 100, Currency: "USD"})

	tReq.Entries[0].Amount.Value = "9223372036854775807"

	res4 := createTransaction(t, &s, tReq)
	defer res4.Body.Close()

	err = json.NewDecoder(res4.Body).Decode(&tRes)
	s.Is.NoErr(err)

	s.Is.True(tRes.Entries[0].Amount == amount.Amount{Value: 9223372036854775807, Currency: "USD"})
}

func TestGetTransaction(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	aReq := accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "cash",
			Type:  account.Asset,
			Basis: "debit",
		},
	}
	res := createAccount(t, &s, aReq)
	defer res.Body.Close()
	var aRes accountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.Is.NoErr(err)
	cashAccount := aRes.Account

	aReq = accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "revenue",
			Type:  account.Asset,
			Basis: "credit",
		},
	}
	res2 := createAccount(t, &s, aReq)
	defer res2.Body.Close()
	err = json.NewDecoder(res2.Body).Decode(&aRes)
	s.Is.NoErr(err)
	revenueAccount := aRes.Account

	tReq := transactionRequest{
		Transaction: api.Transaction{
			Description: "test",
			Entries: []api.Entry{
				{
					Description:   "",
					DebitAccount:  cashAccount.ID,
					CreditAccount: revenueAccount.ID,
					Amount: api.Amount{
						Value:    "100",
						Currency: "USD",
					},
				},
			},
		},
	}

	res3 := createTransaction(t, &s, tReq)
	defer res3.Body.Close()

	var tRes transactionResponse
	err = json.NewDecoder(res3.Body).Decode(&tRes)
	s.Is.NoErr(err)

	s.Is.True(tRes != (transactionResponse{}))

	s.Is.True(tRes.Entries[0].CreditAccount == revenueAccount.ID)
	s.Is.True(tRes.Entries[0].DebitAccount == cashAccount.ID)
	s.Is.True(tRes.Entries[0].Amount == amount.Amount{Value: 100, Currency: "USD"})

	tReq.Entries[0].Amount.Value = "9223372036854775807"
	tReq.Entries[0].Amount.Currency = "usd"

	res4 := getTransfer(t, &s, tRes.ID)
	defer res4.Body.Close()
	s.Is.True(res4.StatusCode == http.StatusOK)

	var getTRes transactionResponse
	err = json.NewDecoder(res4.Body).Decode(&getTRes)
	s.Is.NoErr(err)

	s.Is.Equal(getTRes, tRes)
}

func TestCreateTransaction_int64_overflow(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	aReq := accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "cash",
			Type:  account.Asset,
			Basis: "debit",
		},
	}
	res := createAccount(t, &s, aReq)
	defer res.Body.Close()
	var aRes accountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.Is.NoErr(err)
	cashAccount := aRes.Account

	aReq = accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "revenue",
			Type:  account.Asset,
			Basis: "credit",
		},
	}
	res2 := createAccount(t, &s, aReq)
	defer res2.Body.Close()
	err = json.NewDecoder(res2.Body).Decode(&aRes)
	s.Is.NoErr(err)
	revenueAccount := aRes.Account

	tReq := transactionRequest{
		Transaction: api.Transaction{
			Description: "test",
			Entries: []api.Entry{
				{
					Description:   "",
					DebitAccount:  cashAccount.ID,
					CreditAccount: revenueAccount.ID,
					Amount: api.Amount{
						Value:    "9223372036854775808",
						Currency: "USD",
					},
				},
			},
		},
	}

	res3 := createTransaction(t, &s, tReq)
	defer res3.Body.Close()

	var errRes ErrorResponse
	err = json.NewDecoder(res3.Body).Decode(&errRes)
	s.Is.NoErr(err)

	s.Is.True(errRes != (ErrorResponse{}))
	s.Is.True(strings.Contains(errRes.ValidationErrors, "int64"))
}

func createTransaction(
	t *testing.T,
	s *test.Scope,
	e interface{},
) *http.Response {
	t.Helper()

	var body = new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(e)
	s.Is.NoErr(err)

	req, err := http.NewRequest(
		http.MethodPost,
		s.HTTPTestServer.URL+"/transactions",
		body,
	)
	s.Is.NoErr(err)

	res, err := s.HTTPClient.Do(req)
	s.Is.NoErr(err)

	return res
}

func getTransfer(
	t *testing.T,
	s *test.Scope,
	id string,
) *http.Response {
	t.Helper()

	req, err := http.NewRequest(
		http.MethodGet,
		s.HTTPTestServer.URL+"/transactions/"+id,
		nil,
	)
	s.Is.NoErr(err)

	res, err := s.HTTPClient.Do(req)
	s.Is.NoErr(err)

	return res
}
