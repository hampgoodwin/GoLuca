package api

type Amount struct {
	Value    string `json:"value" validate:"int64,gte=0"`
	Currency string `json:"currency" validate:"len=3,alpha"`
}
