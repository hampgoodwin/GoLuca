package validate

var (
	account = map[string]string{
		"Id":        "required,uuid7",
		"parentId":  "omitempty,uuid7",
		"Name":      "required",
		"Type":      "required",
		"Basis":     "required",
		"CreatedAt": "required",
	}
	getAccountRequest = map[string]string{
		"AccountId": "required,uuid7",
	}
	createAccountRequest = map[string]string{
		"ParentId": "omitempty,uuid7",
		"Name":     "required",
		"Type":     "required",
		"Basis":    "required",
	}
	getTransactionRequest = map[string]string{
		"TransactionId": "required,uuid7",
	}
	transaction = map[string]string{
		"Id":          "required,uuid7",
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
		"DebitAccount":  "required,uuid7",
		"CreditAccount": "required,uuid7",
		"Amount":        "required",
	}
	amount = map[string]string{
		"Value":    "gte=0",
		"Currency": "len=3,alpha",
	}
)
