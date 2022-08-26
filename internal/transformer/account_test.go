package transformer

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/pkg/account"
	httpaccount "github.com/hampgoodwin/GoLuca/pkg/http/v0/account"
	"github.com/matryer/is"
)

func TestNewAccountFromHTTPCreateAccount(t *testing.T) {
	parentID := uuid.NewString()
	testCases := []struct {
		description       string
		httpCreateAccount httpaccount.CreateAccount
		expected          account.Account
	}{
		{description: "empty"},
		{
			description: "success",
			httpCreateAccount: httpaccount.CreateAccount{
				ParentID: parentID,
				Name:     "asset",
				Type:     "asset",
				Basis:    "debit",
			},
			expected: account.Account{
				ParentID: parentID,
				Name:     "asset",
				Type:     account.TypeAsset,
				Basis:    "debit",
			},
		},
	}

	a := is.New(t)

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := NewAccountFromHTTPCreateAccount(tc.httpCreateAccount)

			a.Equal(tc.expected, actual)
		})
	}
}

func TestNewHTTPAccountFromAccount(t *testing.T) {
	testCases := []struct {
		description string
		account     account.Account
		expected    httpaccount.Account
	}{
		{description: "empty"},
		{
			description: "success",
			account: account.Account{
				ID:       "ID",
				ParentID: "parentID",
				Name:     "equity",
				Type:     account.TypeEquity,
				Basis:    "credit",
			},
			expected: httpaccount.Account{
				ID:       "ID",
				ParentID: "parentID",
				Name:     "equity",
				Type:     "equity",
				Basis:    "credit",
			},
		},
	}

	a := is.New(t)

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := NewHTTPAccountFromAccount(tc.account)

			a.Equal(tc.expected, actual)
		})
	}
}
