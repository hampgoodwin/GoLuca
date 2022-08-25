package transformer

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/pkg/account"
	httpaccount "github.com/hampgoodwin/GoLuca/pkg/http/account"
	"github.com/matryer/is"
)

func TestNewAccountFromHTTPAccount(t *testing.T) {
	parentID := uuid.NewString()
	testCases := []struct {
		description       string
		httpCreateAccount httpaccount.CreateAccount
		expected          account.Account
		err               error
	}{
		{description: "empty"},
		{
			description: "success",
			httpCreateAccount: httpaccount.CreateAccount{
				ParentID: parentID,
				Name:     "asset",
				Type:     account.Asset,
				Basis:    "debit",
			},
			expected: account.Account{
				ParentID: parentID,
				Name:     "asset",
				Type:     account.Asset,
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
