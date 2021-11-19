package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/internal/controller"
	"github.com/hampgoodwin/GoLuca/internal/test"
	"github.com/hampgoodwin/GoLuca/pkg/account"
)

func TestCreateAccount(t *testing.T) {
	s := test.GetScope(t)

	aReq := controller.AccountRequest{
		Account: &account.Account{
			Name:  "cash",
			Type:  account.Asset,
			Basis: "debit",
		},
	}

	res := createAccount(t, &s, aReq)
	defer res.Body.Close()

	var aRes controller.AccountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.IS.NoErr(err)

	s.IS.True(aRes != (controller.AccountResponse{}))

	s.IS.Equal(aReq.Account.Name, aRes.Account.Name)
	s.IS.Equal(aReq.Account.Type, aRes.Account.Type)
	s.IS.Equal(aReq.Account.Basis, aRes.Account.Basis)
	s.IS.True(aRes.Account.ID != "")
	s.IS.True(aRes.Account.ParentID == "")
}

func TestCreateAccount_InvalidRequestBody(t *testing.T) {
	s := test.GetScope(t)

	aReq := controller.AccountRequest{
		Account: &account.Account{
			Name:  "",
			Type:  account.Type("type"),
			Basis: "sandwhich",
		},
	}

	res := createAccount(t, &s, aReq)
	defer res.Body.Close()

	var errRes controller.ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errRes)
	s.IS.NoErr(err)

	s.IS.Equal("validating deserialized account body", errRes.Description)
	s.IS.Equal("Key: 'Account.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'Account.Type' Error:Field validation for 'Type' failed on the 'oneof' tag\nKey: 'Account.Basis' Error:Field validation for 'Basis' failed on the 'oneof' tag", errRes.ValidationErrors)
}

func TestCreateAccount_CannotDeserialize(t *testing.T) {
	s := test.GetScope(t)

	bad := []byte{}

	res := createAccount(t, &s, bad)
	defer res.Body.Close()

	var errRes controller.ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errRes)
	s.IS.NoErr(err)

	s.IS.Equal("json: cannot unmarshal string into Go value of type controller.AccountRequest", errRes.Description)

}

func TestGetAccount(t *testing.T) {
	s := test.GetScope(t)

	// Create an account and assert
	aReq := controller.AccountRequest{
		Account: &account.Account{
			Name:  "cash",
			Type:  account.Asset,
			Basis: "debit",
		},
	}

	res := createAccount(t, &s, aReq)
	defer res.Body.Close()

	var aRes controller.AccountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.IS.NoErr(err)

	// Get the created account
	getRes := getAccount(t, &s, aRes.ID)
	res.Body.Close()

	var getARes controller.AccountResponse
	err = json.NewDecoder(getRes.Body).Decode(&getARes)
	s.IS.NoErr(err)

	s.IS.Equal(aRes, getARes)
}

func TestGetAccount_ErrorNotFound(t *testing.T) {
	s := test.GetScope(t)

	id := uuid.NewString()
	res := getAccount(t, &s, id)

	var errRes controller.ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errRes)
	s.IS.NoErr(err)
	s.IS.Equal(fmt.Sprintf("account %q not found", id), errRes.Description)
}

func TestGetAccount_InvalidPersistedAccount(t *testing.T) {
	s := test.GetScope(t)

	id := gofakeit.Sentence(1)
	parentID := gofakeit.Sentence(1)
	name := gofakeit.Name()
	typ := gofakeit.Sentence(1)
	basis := strings.Replace(gofakeit.Sentence(5), " ", "", -1)[0:5]

	_, err := s.DB.Exec(s.CTX, `
	INSERT INTO account (id, parent_id, name, type, basis)
		VALUES($1, $2, $3, $4, $5)
	;`, id, parentID, name, typ, basis)
	s.IS.NoErr(err)

	res := getAccount(t, &s, id)
	defer res.Body.Close()

	s.IS.Equal(http.StatusNotFound, res.StatusCode)
}

func TestGetAccounts(t *testing.T) {
	s := test.GetScope(t)

	// Create an account and assert
	aReq := controller.AccountRequest{
		Account: &account.Account{
			Name:  "cash",
			Type:  account.Asset,
			Basis: "debit",
		},
	}
	res := createAccount(t, &s, aReq)
	defer res.Body.Close()

	var a1 controller.AccountResponse
	err := json.NewDecoder(res.Body).Decode(&a1)
	s.IS.NoErr(err)

	aReq.Name = "accounts receivable"
	res2 := createAccount(t, &s, aReq)
	defer res2.Body.Close()

	var a2 controller.AccountResponse
	err = json.NewDecoder(res2.Body).Decode(&a2)
	s.IS.NoErr(err)

	aRes := getAccounts(t, &s)
	s.IS.True(len(aRes.Accounts) == 2)

	i := 0
	for _, a := range aRes.Accounts {
		if a == *a1.Account || a == *a2.Account {
			i++
		}
	}
	s.IS.True(i == len(aRes.Accounts))
}

func createAccount(
	t *testing.T,
	s *test.Scope,
	e interface{},
) *http.Response {
	t.Helper()

	var body = new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(e)
	s.IS.NoErr(err)

	req, err := http.NewRequest(
		http.MethodPost,
		"http://"+s.Env.Config.HTTPAPI.AddressString()+"/accounts",
		body,
	)
	s.IS.NoErr(err)

	res, err := s.HTTPClient.Do(req)
	s.IS.NoErr(err)

	return res
}

func getAccount(
	t *testing.T,
	s *test.Scope,
	id string,
) *http.Response {
	// Get the created account and assert it's equal to the created account
	req, err := http.NewRequest(
		http.MethodGet,
		"http://"+s.Env.Config.HTTPAPI.AddressString()+"/accounts/"+id,
		nil,
	)
	s.IS.NoErr(err)

	res, err := s.HTTPClient.Do(req)
	s.IS.NoErr(err)

	return res
}

func getAccounts(
	t *testing.T,
	s *test.Scope,
) controller.AccountsResponse {
	// Get the created account and assert it's equal to the created account
	req, err := http.NewRequest(
		http.MethodGet,
		"http://"+s.Env.Config.HTTPAPI.AddressString()+"/accounts",
		nil,
	)
	s.IS.NoErr(err)

	res, err := s.HTTPClient.Do(req)
	s.IS.NoErr(err)

	var aRes controller.AccountsResponse
	err = json.NewDecoder(res.Body).Decode(&aRes)
	s.IS.NoErr(err)

	return aRes
}
