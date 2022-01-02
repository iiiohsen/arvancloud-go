package sdk

import "context"

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

// GetAccount gets the contact and billing information related to the Account.
func (c *Client) GetAccount(ctx context.Context) (*Account, error) {
	e, err := c.Account.Endpoint()
	if err != nil {
		return nil, err
	}

	r, err := coupleAPIErrors(c.R(ctx).SetResult(&Account{}).Get(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*Account), nil
}
