package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

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

	var body = new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&aReq)
	s.IS.NoErr(err)

	req, err := http.NewRequest(
		http.MethodPost,
		"http://"+s.Env.Config.HTTPAPI.AddressString()+"/accounts",
		body,
	)
	s.IS.NoErr(err)

	res, err := s.HTTPClient.Do(req)
	s.IS.NoErr(err)

	var aRes controller.AccountResponse
	err = json.NewDecoder(res.Body).Decode(&aRes)
	s.IS.NoErr(err)
	res.Body.Close()

	s.IS.Equal(aReq.Account.Name, aRes.Account.Name)
	s.IS.Equal(aReq.Account.Type, aRes.Account.Type)
	s.IS.Equal(aReq.Account.Basis, aRes.Account.Basis)
	s.IS.True(aRes.Account.ID != "")
	s.IS.True(aRes.Account.ParentID == "")
}
