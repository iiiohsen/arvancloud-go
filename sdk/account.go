package sdk

// Account associated with the token in use.
type Account struct {
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	Email             string  `json:"email"`
	Company           string  `json:"company"`
	Address1          string  `json:"address_1"`
	Address2          string  `json:"address_2"`
	Balance           float32 `json:"balance"`
	BalanceUninvoiced float32 `json:"balance_uninvoiced"`
	City              string  `json:"city"`
	State             string  `json:"state"`
	Zip               string  `json:"zip"`
	Country           string  `json:"country"`
	TaxID             string  `json:"tax_id"`
	Phone             string  `json:"phone"`
}
