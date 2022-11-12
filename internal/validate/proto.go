package validate

var (
	listAccountsResponseValidation = map[string]string{
		"Accounts": "required",
	}
	accountValidation = map[string]string{
		"Id":        "required,KSUID",
		"parentId":  "omitempty,KSUID",
		"Name":      "required",
		"Type":      "required",
		"Basis":     "required",
		"CreatedAt": "required",
	}
	getAccountRequest = map[string]string{
		"AccountId": "required,KSUID",
	}
	createAccountRequest = map[string]string{
		"ParentId": "omitempty,KSUID",
		"Name":     "required",
		"Type":     "required",
		"Basis":    "required",
	}
)
