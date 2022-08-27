package amount

type Amount struct {
	Value    string `json:"value" validate:"required,stringAsInt64"`
	Currency string `json:"currency" validate:"required,len=3"`
}
