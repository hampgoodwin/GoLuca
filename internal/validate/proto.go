package validate

var (
	account = map[string]string{
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
	getTransactionRequest = map[string]string{
		"TransactionId": "required,KSUID",
	}
	transaction = map[string]string{
		"Id":          "required,KSUID",
		"Description": "required",
		"Entries":     "required,gt=0,dive",
		"CreatedAt":   "required",
	}
	createTransactionRequest = map[string]string{
		"Description": "required",
		"Entries":     "required,gt=0,dive",
	}
	createEntry = map[string]string{
		"Description":   "required",
		"DebitAccount":  "required,KSUID",
		"CreditAccount": "required,KSUID",
		"Amount":        "required",
	}
	amount = map[string]string{
		"Value":    "gte=0",
		"Currency": "len=3,alpha",
	}
)
