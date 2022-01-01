package arvancloud

import (
	"log"
	"net/http"
	"os"

	"github.com/S4eedb/arvancloud-go/sdk"

	"golang.org/x/oauth2"
)

func New(token, email string) (api sdk.Client, err error) {
	if token == "" || email == "" {
		err = sdk.NewError("empty credential")
		return
	}
	apiKey, ok := os.LookupEnv("ARVANCLOUD_TOKEN")
	if !ok {
		log.Fatal("Could not find ARVANCLOUD_TOKEN, please assert it is set.")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}
	NewApi := sdk.NewClient(oauth2Client)

	// api.APIKey = token
	// api.APIEmail = email
	// api.authType = AuthKeyEmail

	return NewApi, nil
}
