package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hampgoodwin/GoLuca/internal/httpapi"
	"github.com/hampgoodwin/GoLuca/internal/test"
	"github.com/hampgoodwin/GoLuca/pkg/account"
	"github.com/hampgoodwin/GoLuca/pkg/amount"
)

func TestCreateTransaction(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	aReq := accountRequest{
		Account: &account.Account{
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
		Account: &account.Account{
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
		Transaction: httpapi.Transaction{
			Description: "test",
			Entries: []httpapi.Entry{
				{
					Description:   "",
					DebitAccount:  cashAccount.ID,
					CreditAccount: revenueAccount.ID,
					Amount: httpapi.Amount{
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
