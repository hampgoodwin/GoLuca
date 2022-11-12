package validate

var (
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
