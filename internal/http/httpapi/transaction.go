package api

type Transaction struct {
	Description string  `json:"description" validate:"required"`
	Entries     []Entry `json:"entries,omitempty" validate:"dive,gte=1"`
}

func (t Transaction) IsZero() bool {
	if t.Description != "" {
		return false
	}
	if t.Entries != nil {
		return false
	}
	return true
}
