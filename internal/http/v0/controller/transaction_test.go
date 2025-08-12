package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hampgoodwin/GoLuca/internal/test"
	httpaccount "github.com/hampgoodwin/GoLuca/pkg/http/v0/account"
	httpamount "github.com/hampgoodwin/GoLuca/pkg/http/v0/amount"
	httptransaction "github.com/hampgoodwin/GoLuca/pkg/http/v0/transaction"
)

func TestCreateTransaction(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	aReq := accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "cash",
			Type:  "asset",
			Basis: "debit",
		},
	}
	res := createAccount(t, &s, aReq)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()
	var aRes accountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.Is.NoErr(err)
	cashAccount := aRes.Account

	aReq = accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "revenue",
			Type:  "asset",
			Basis: "credit",
		},
	}
	res2 := createAccount(t, &s, aReq)
	defer func() {
		if err := res2.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()
	err = json.NewDecoder(res2.Body).Decode(&aRes)
	s.Is.NoErr(err)
	revenueAccount := aRes.Account

	tReq := transactionRequest{
		Transaction: httptransaction.CreateTransaction{
			Description: "test",
			Entries: []httptransaction.CreateEntry{
				{
					Description:   "",
					DebitAccount:  cashAccount.ID,
					CreditAccount: revenueAccount.ID,
					Amount: httpamount.Amount{
						Value:    "100",
						Currency: "USD",
					},
				},
			},
		},
	}

	res3 := createTransaction(t, &s, tReq)
	defer func() {
		if err := res3.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()

	var tRes transactionResponse
	err = json.NewDecoder(res3.Body).Decode(&tRes)
	s.Is.NoErr(err)

	// s.Is.True(tRes != (transactionResponse{}))

	s.Is.True(tRes.Entries[0].CreditAccount == revenueAccount.ID)
	s.Is.True(tRes.Entries[0].DebitAccount == cashAccount.ID)
	s.Is.True(tRes.Entries[0].Amount == httpamount.Amount{Value: "100", Currency: "USD"})

	tReq.Transaction.Entries[0].Amount.Value = "9223372036854775807"

	res4 := createTransaction(t, &s, tReq)
	defer func() {
		if err := res4.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()

	err = json.NewDecoder(res4.Body).Decode(&tRes)
	s.Is.NoErr(err)

	s.Is.True(tRes.Entries[0].Amount == httpamount.Amount{Value: "9223372036854775807", Currency: "USD"})
}

func TestCreateTransaction_int64_overflow(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	aReq := accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "cash",
			Type:  "asset",
			Basis: "debit",
		},
	}
	res := createAccount(t, &s, aReq)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()
	var aRes accountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.Is.NoErr(err)
	cashAccount := aRes.Account

	aReq = accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "revenue",
			Type:  "asset",
			Basis: "credit",
		},
	}
	res2 := createAccount(t, &s, aReq)
	defer func() {
		if err := res2.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()
	err = json.NewDecoder(res2.Body).Decode(&aRes)
	s.Is.NoErr(err)
	revenueAccount := aRes.Account

	tReq := transactionRequest{
		Transaction: httptransaction.CreateTransaction{
			Description: "test",
			Entries: []httptransaction.CreateEntry{
				{
					Description:   "",
					DebitAccount:  cashAccount.ID,
					CreditAccount: revenueAccount.ID,
					Amount: httpamount.Amount{
						Value:    "9223372036854775808",
						Currency: "USD",
					},
				},
			},
		},
	}

	res3 := createTransaction(t, &s, tReq)
	defer func() {
		if err := res3.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()

	var errRes ErrorResponse
	err = json.NewDecoder(res3.Body).Decode(&errRes)
	s.Is.NoErr(err)

	s.Is.True(errRes != (ErrorResponse{}))
	fmt.Println(errRes)
	s.Is.True(strings.Contains(errRes.ValidationErrors, "stringAsInt64"))
}

func TestGetTransaction(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	aReq := accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "cash",
			Type:  "asset",
			Basis: "debit",
		},
	}
	res := createAccount(t, &s, aReq)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()
	var aRes accountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.Is.NoErr(err)
	cashAccount := aRes.Account

	aReq = accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "revenue",
			Type:  "asset",
			Basis: "credit",
		},
	}
	res2 := createAccount(t, &s, aReq)
	defer func() {
		if err := res2.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()
	err = json.NewDecoder(res2.Body).Decode(&aRes)
	s.Is.NoErr(err)
	revenueAccount := aRes.Account

	tReq := transactionRequest{
		Transaction: httptransaction.CreateTransaction{
			Description: "test",
			Entries: []httptransaction.CreateEntry{
				{
					Description:   "",
					DebitAccount:  cashAccount.ID,
					CreditAccount: revenueAccount.ID,
					Amount: httpamount.Amount{
						Value:    "100",
						Currency: "USD",
					},
				},
			},
		},
	}

	res3 := createTransaction(t, &s, tReq)
	defer func() {
		if err := res3.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()

	var tRes transactionResponse
	err = json.NewDecoder(res3.Body).Decode(&tRes)
	s.Is.NoErr(err)

	// s.Is.True(tRes != (transactionResponse{}))

	s.Is.True(tRes.Entries[0].CreditAccount == revenueAccount.ID)
	s.Is.True(tRes.Entries[0].DebitAccount == cashAccount.ID)
	s.Is.True(tRes.Entries[0].Amount == httpamount.Amount{Value: "100", Currency: "USD"})

	tReq.Transaction.Entries[0].Amount.Value = "9223372036854775807"
	tReq.Transaction.Entries[0].Amount.Currency = "usd"

	res4 := getTransaction(t, &s, tRes.ID)
	defer func() {
		if err := res4.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()
	s.Is.True(res4.StatusCode == http.StatusOK)

	var getTRes transactionResponse
	err = json.NewDecoder(res4.Body).Decode(&getTRes)
	s.Is.NoErr(err)

	s.Is.Equal(getTRes, tRes)
}

func TestListTransactions(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	tsRes := transactionsResponse{
		Transactions: []httptransaction.Transaction{},
	}

	aReq := accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "cash",
			Type:  "asset",
			Basis: "debit",
		},
	}
	res := createAccount(t, &s, aReq)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()
	var aRes accountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.Is.NoErr(err)
	cashAccount := aRes.Account

	aReq = accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "revenue",
			Type:  "asset",
			Basis: "credit",
		},
	}
	res2 := createAccount(t, &s, aReq)
	defer func() {
		if err := res2.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()
	err = json.NewDecoder(res2.Body).Decode(&aRes)
	s.Is.NoErr(err)
	revenueAccount := aRes.Account

	tReq := transactionRequest{
		Transaction: httptransaction.CreateTransaction{
			Description: "test",
			Entries: []httptransaction.CreateEntry{
				{
					Description:   "",
					DebitAccount:  cashAccount.ID,
					CreditAccount: revenueAccount.ID,
					Amount: httpamount.Amount{
						Value:    "100",
						Currency: "USD",
					},
				},
			},
		},
	}

	res3 := createTransaction(t, &s, tReq)
	defer func() {
		if err := res3.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()

	var tRes transactionResponse
	err = json.NewDecoder(res3.Body).Decode(&tRes)
	s.Is.NoErr(err)
	tsRes.Transactions = append(tsRes.Transactions, tRes.Transaction)

	// s.Is.True(tRes != (transactionResponse{}))

	s.Is.True(tRes.Entries[0].CreditAccount == revenueAccount.ID)
	s.Is.True(tRes.Entries[0].DebitAccount == cashAccount.ID)
	s.Is.True(tRes.Entries[0].Amount == httpamount.Amount{Value: "100", Currency: "USD"})

	tReq.Transaction.Entries[0].Amount.Value = "9223372036854775807"
	tReq.Transaction.Entries[0].Amount.Currency = "usd"

	res4 := listTransactions(t, &s)
	defer func() {
		if err := res4.Body.Close(); err != nil {
			log.Printf("creating account: %v\n", err)
		}
	}()
	s.Is.True(res4.StatusCode == http.StatusOK)

	var getTsRes transactionsResponse
	err = json.NewDecoder(res4.Body).Decode(&getTsRes)
	s.Is.NoErr(err)

	s.Is.Equal(getTsRes, tsRes)
}

func createTransaction(
	t *testing.T,
	s *test.Scope,
	e any,
) *http.Response {
	t.Helper()

	body := new(bytes.Buffer)
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

func getTransaction(
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

func listTransactions(
	t *testing.T,
	s *test.Scope,
) *http.Response {
	t.Helper()

	req, err := http.NewRequest(
		http.MethodGet,
		s.HTTPTestServer.URL+"/transactions/",
		nil,
	)
	s.Is.NoErr(err)

	res, err := s.HTTPClient.Do(req)
	s.Is.NoErr(err)

	return res
}
