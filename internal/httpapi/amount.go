package httpapi

type Amount struct {
	Value    string `json:"value" validate:"gte=0"`
	Currency string `json:"currency" validate:"len=3,alpha"`
}
