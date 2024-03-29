package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/segmentio/ksuid"

	"github.com/hampgoodwin/GoLuca/internal/test"
	httpaccount "github.com/hampgoodwin/GoLuca/pkg/http/v0/account"
)

func TestCreateAccount(t *testing.T) {
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
	defer res.Body.Close()

	var aRes accountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.Is.NoErr(err)

	s.Is.True(aRes != (accountResponse{}))

	s.Is.Equal(aReq.Account.Name, aRes.Account.Name)
	s.Is.Equal(aReq.Account.Type, aRes.Account.Type)
	s.Is.Equal(aReq.Account.Basis, aRes.Account.Basis)
	s.Is.True(aRes.Account.ID != "")
	s.Is.True(aRes.Account.ParentID == "")
}

func TestCreateAccount_InvalidRequestBody(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	aReq := accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "",
			Type:  "",
			Basis: "sandwhich",
		},
	}

	res := createAccount(t, &s, aReq)
	defer res.Body.Close()

	var errRes ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errRes)
	s.Is.NoErr(err)

	s.Is.Equal("validating account", errRes.Description)
	s.Is.Equal("Key: 'Account.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'Account.Type.Slug' Error:Field validation for 'Slug' failed on the 'oneof' tag\nKey: 'Account.Basis.Slug' Error:Field validation for 'Slug' failed on the 'oneof' tag", errRes.ValidationErrors)
}

func TestCreateAccount_CannotDeserialize(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	bad := []byte{}

	res := createAccount(t, &s, bad)
	defer res.Body.Close()

	var errRes ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errRes)
	s.Is.NoErr(err)

	s.Is.Equal("json: cannot unmarshal string into Go value of type controller.accountRequest", errRes.Description)

}

func TestGetAccount(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	// Create an account and assert
	aReq := accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "cash",
			Type:  "asset",
			Basis: "debit",
		},
	}

	res := createAccount(t, &s, aReq)
	defer res.Body.Close()

	var aRes accountResponse
	err := json.NewDecoder(res.Body).Decode(&aRes)
	s.Is.NoErr(err)

	// Get the created account
	getRes := getAccount(t, &s, aRes.ID)
	res.Body.Close()

	var getARes accountResponse
	err = json.NewDecoder(getRes.Body).Decode(&getARes)
	s.Is.NoErr(err)

	s.Is.Equal(aRes, getARes)
}

func TestGetAccount_ErrorNotFound(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	id := ksuid.New().String()
	res := getAccount(t, &s, id)

	var errRes ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errRes)
	s.Is.NoErr(err)
	s.Is.Equal(fmt.Sprintf("account %q not found", id), errRes.Description)
}

func TestGetAccount_InvalidPersistedAccount(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	id := gofakeit.Sentence(1)
	parentID := gofakeit.Sentence(1)
	name := gofakeit.Name()
	typ := gofakeit.Sentence(1)
	basis := strings.Replace(gofakeit.Sentence(5), " ", "", -1)[0:5]

	_, err := s.DB.Exec(s.Ctx, `
	INSERT INTO account (id, parent_id, name, type, basis)
		VALUES($1, $2, $3, $4, $5)
	;`, id, parentID, name, typ, basis)
	s.Is.NoErr(err)

	res := getAccount(t, &s, id)
	defer res.Body.Close()

	s.Is.Equal(http.StatusNotFound, res.StatusCode)
}

func TestListAccounts(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	// Create an account and assert
	aReq := accountRequest{
		Account: httpaccount.CreateAccount{
			Name:  "cash",
			Type:  "asset",
			Basis: "debit",
		},
	}
	res := createAccount(t, &s, aReq)
	defer res.Body.Close()

	var a1 accountResponse
	err := json.NewDecoder(res.Body).Decode(&a1)
	s.Is.NoErr(err)

	aReq.Account.Name = "accounts receivable"
	res2 := createAccount(t, &s, aReq)
	defer res2.Body.Close()

	var a2 accountResponse
	err = json.NewDecoder(res2.Body).Decode(&a2)
	s.Is.NoErr(err)

	httpResponse := listAccounts(t, &s, "", "")

	var aRes accountsResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&aRes)
	s.Is.NoErr(err)

	s.Is.True(len(aRes.Accounts) == 2)

	i := 0
	for _, a := range aRes.Accounts {
		if a == a1.Account || a == a2.Account {
			i++
		}
	}
	s.Is.True(i == len(aRes.Accounts))
}

func TestListAccounts_InvalidRequestBody(t *testing.T) {
	s := test.GetScope(t)
	s.SetHTTP(t, newTestHTTPHandler(s.Env.Log, s.DB))

	httpResponse := listAccounts(t, &s, "invalid_cursor", "")

	var errorResponse ErrorResponse
	err := json.NewDecoder(httpResponse.Body).Decode(&errorResponse)
	s.Is.NoErr(err)

	s.Is.True(errorResponse.Description == "invalid cursor or token")
	s.Is.True(httpResponse.StatusCode == http.StatusBadRequest)
}

func createAccount(
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
		s.HTTPTestServer.URL+"/accounts",
		body,
	)
	s.Is.NoErr(err)

	res, err := s.HTTPClient.Do(req)
	s.Is.NoErr(err)

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
		s.HTTPTestServer.URL+"/accounts/"+id,
		nil,
	)
	s.Is.NoErr(err)

	res, err := s.HTTPClient.Do(req)
	s.Is.NoErr(err)

	return res
}

func listAccounts(
	t *testing.T,
	s *test.Scope,
	cursor string,
	limit string,
) *http.Response {
	// Get the created account and assert it's equal to the created account
	req, err := http.NewRequest(
		http.MethodGet,
		s.HTTPTestServer.URL+"/accounts"+fmt.Sprintf("?cursor=%s&limit=%s", cursor, limit),
		nil,
	)
	s.Is.NoErr(err)

	res, err := s.HTTPClient.Do(req)
	s.Is.NoErr(err)

	return res
}
