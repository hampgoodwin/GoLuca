package api

type Entry struct {
	Description   string `json:"description"`
	DebitAccount  string `json:"debitAccount" validate:"required,uuid4"`
	CreditAccount string `json:"creditAccount" validate:"required,uuid4"`
	Amount        Amount `json:"amount" validate:"required"`
}
