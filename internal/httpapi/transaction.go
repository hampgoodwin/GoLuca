package httpapi

type Transaction struct {
	Description string  `json:"description" validate:"required"`
	Entries     []Entry `json:"entries,omitempty" validate:"dive,gte=1"`
}
