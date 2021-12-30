package sdk

import "errors"

const apiURL = "https://napi.arvancloud.com"

// API holds the configuration for the current API client. A client should not
// be modified concurrently.
type API struct {
	APIEmail string
	APIToken string
	BaseURL  string
}

// newClient provides shared logic for New and NewWithUserServiceKey
func newClient() (*API, error) {

	api := &API{
		BaseURL: apiURL,
	}

	return api, nil
}

// New creates a new Cloudflare v4 API client.
func New(token, email string) (*API, error) {
	if token == "" || email == "" {
		return nil, errors.New("err Empty Credentials")
	}

	api, err := newClient()
	if err != nil {
		return nil, err
	}

	api.APIToken = token
	api.APIEmail = email

	return api, nil
}
