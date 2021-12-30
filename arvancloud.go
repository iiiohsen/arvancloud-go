package arvancloud

import (
	"fmt"

	"github.com/S4eedb/arvancloud-go/sdk/client"
	"github.com/S4eedb/arvancloud-go/sdk/errors"
)

func New(token, email string) (api string, err error) {
	if token == "" || email == "" {
		EmptyCredentials := fmt.Sprintf(errors.EmptyCredentialsMessage)
		err = errors.NewClientError(errors.EmptyCredentialsCode, EmptyCredentials, nil)
		return
	}

	NewApi, err := client.NewClient()
	if err != nil {
		return "", err
	}

	// api.APIKey = token
	// api.APIEmail = email
	// api.authType = AuthKeyEmail

	return NewApi, nil
}
