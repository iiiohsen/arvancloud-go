package sdk

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/go-resty/resty/v2"
)

const (
	accountName         = "account"
	accountSettingsName = "accountsettings"
	domainRecordsName   = "records"
	domainsName         = "domains"
	accountEndpoint     = "account"
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

func (r Resource) render(data ...interface{}) (string, error) {
	if data == nil {
		return "", NewError("Cannot template endpoint with <nil> data")
	}
	out := ""
	buf := bytes.NewBufferString(out)

	var substitutions interface{}

	switch len(data) {
	case 1:
		substitutions = struct{ ID interface{} }{data[0]}
	case 2:
		substitutions = struct {
			ID       interface{}
			SecondID interface{}
		}{data[0], data[1]}
	default:
		return "", NewError("Too many arguments to render template (expected 1 or 2)")
	}

	if err := r.endpointTemplate.Execute(buf, substitutions); err != nil {
		return "", NewError(err)
	}
	return buf.String(), nil
}

// endpointWithParams will return the rendered endpoint string for the resource with provided parameters
func (r Resource) endpointWithParams(params ...interface{}) (string, error) {
	if !r.isTemplate {
		return r.endpoint, nil
	}
	return r.render(params...)
}

// Endpoint will return the non-templated endpoint string for resource
func (r Resource) Endpoint() (string, error) {
	if r.isTemplate {
		return "", NewError(fmt.Sprintf("Tried to get endpoint for %s without providing data for template", r.name))
	}
	return r.endpoint, nil
}
