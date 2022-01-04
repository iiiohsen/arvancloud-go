package sdk

import (
	"testing"
)

func TestResourceEndpoint(t *testing.T) {
	apiKey := "MYFAKEAPIKEY"

	client := NewClient(apiKey)

	r := client.Resource("domains")
	e, err := r.Endpoint()
	if err != nil {
		t.Error("Got error when querying for domains endpoint")
	}

	if e != domainsEndpoint {
		t.Errorf("domains endpoint did not match '%s'", domainsEndpoint)
	}
}
