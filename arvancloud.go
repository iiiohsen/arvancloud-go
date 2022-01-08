package arvancloud

import (
	"github.com/S4eedb/arvancloud-go/sdk"
)

func New(apiKey string) (api sdk.Client, err error) {

	NewApi := sdk.NewClient(apiKey)

	return NewApi, nil
}
