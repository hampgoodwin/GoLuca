package amount

type Amount struct {
	// Value denotes the number of currency in an amount. The last two digits
	// are the cents of the given currency. While this value is of type
	// the limitation is 9223372036854775807 due to database type limitations
	Value    int64  `json:"value" validate:"gte=0"`
	Currency string `json:"currency" validate:"len=3,alpha"`
}
