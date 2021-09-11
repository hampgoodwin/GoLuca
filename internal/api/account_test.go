package api

// func TestAccount_Get(t *testing.T) {
// 	// a := assert.New(t)

// 	testCases := []struct {
// 		description string
// 		create      *accountRequest
// 		expected    *accountResponse
// 		statusCode  int
// 		error       errorResponse
// 	}{
// 		{
// 			description: "success-create-single-and-get",
// 			create: &accountRequest{
// 				Account: &account.Account{
// 					Name:  "account receivable",
// 					Type:  account.Asset,
// 					Basis: "debit",
// 				},
// 			},
// 			expected: &accountResponse{
// 				Account: &account.Account{
// 					ParentID: 0,
// 					Name:     "account receivable",
// 					Type:     account.Asset,
// 					Basis:    "debit",
// 				},
// 			},
// 		},
// 	}

// for _, tc := range testCases {
// 	url := fmt.Sprintf("%s:%s/accounts", "config.Env.APIHost", "config.Env.APIPort")
// 	// req, err := http.Post(url, "application/json", json.NewEncoder())
// }
// }
