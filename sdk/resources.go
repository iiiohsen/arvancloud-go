package sdk

import (
	"context"
	"fmt"
	"text/template"

	"github.com/go-resty/resty/v2"
)

const (
	accountName           = "account"
	accountSettingsName   = "accountsettings"
	domainRecordsName     = "records"
	domainsName           = "domains"
	domainsEndpoint       = "domains"
	accountEndpoint       = "account"
	domainRecordsEndpoint = "domains/{{ .ID }}/records"
)

// Resource represents a arvancloud API resource
type Resource struct {
	name             string
	endpoint         string
	isTemplate       bool
	endpointTemplate *template.Template
	R                func(ctx context.Context) *resty.Request
	PR               func(ctx context.Context) *resty.Request
}

// NewResource is the factory to create a new Resource struct. If it has a template string the useTemplate bool must be set.
func NewResource(client *Client, name string, endpoint string, useTemplate bool, singleType interface{}, pagedType interface{}) *Resource {
	var tmpl *template.Template

	if useTemplate {
		tmpl = template.Must(template.New(name).Parse(endpoint))
	}

	r := func(ctx context.Context) *resty.Request {
		return client.R(ctx).SetResult(singleType)
	}

	pr := func(ctx context.Context) *resty.Request {
		return client.R(ctx).SetResult(pagedType)
	}

	return &Resource{name, endpoint, useTemplate, tmpl, r, pr}
}

// Endpoint will return the non-templated endpoint string for resource
func (r Resource) Endpoint() (string, error) {
	if r.isTemplate {
		return "", NewError(fmt.Sprintf("Tried to get endpoint for %s without providing data for template", r.name))
	}
	return r.endpoint, nil
}
