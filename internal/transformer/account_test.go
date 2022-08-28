package transformer

import (
	"fmt"
	"testing"

	"github.com/hampgoodwin/GoLuca/internal/account"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	httpaccount "github.com/hampgoodwin/GoLuca/pkg/http/v0/account"
	"github.com/matryer/is"
	"github.com/segmentio/ksuid"
)

func TestNewAccountFromHTTPCreateAccount(t *testing.T) {
	parentID := ksuid.New().String()
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
				Basis:    account.BasisDebit,
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
				Basis:    account.BasisCredit,
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

func TestNewAccountFromRepoAccount(t *testing.T) {
	testCases := []struct {
		description string
		account     repository.Account
		expected    account.Account
	}{
		{description: "empty"},
		{
			description: "success",
			account: repository.Account{
				ID:       "ID",
				ParentID: "parentID",
				Name:     "equity",
				Type:     "equity",
				Basis:    "credit",
			},
			expected: account.Account{
				ID:       "ID",
				ParentID: "parentID",
				Name:     "equity",
				Type:     account.TypeEquity,
				Basis:    account.BasisCredit,
			},
		},
	}

	a := is.New(t)

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			t.Parallel()
			actual := NewAccountFromRepoAccount(tc.account)

			a.Equal(tc.expected, actual)
		})
	}
}

func TestNewRepoAccountFromAccount(t *testing.T) {
	testCases := []struct {
		description string
		account     account.Account
		expected    repository.Account
	}{
		{description: "empty"},
		{
			description: "success",
			account: account.Account{
				ID:       "ID",
				ParentID: "parentID",
				Name:     "equity",
				Type:     account.TypeEquity,
				Basis:    account.BasisCredit,
			},
			expected: repository.Account{
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
			actual := NewRepoAccountFromAccount(tc.account)

			a.Equal(tc.expected, actual)
		})
	}
}
