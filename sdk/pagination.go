package sdk

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// PageOptions are the pagination parameters for List endpoints
type PageOptions struct {
	CurrentPage int    `json:"current_page"`
	From        int    `json:"from"`
	LastPage    int    `json:"last_page"`
	Path        string `json:"path"`
	PageSize    int    `json:"per_page"`
	To          int    `json:"to"`
	Total       int    `json:"total"`
}

// ListOptions are the pagination and filtering (TODO) parameters for endpoints
type ListOptions struct {
	*PageOptions
}

// NewListOptions simplified construction of ListOptions using only
// the two writable properties, Page and Filter
func NewListOptions(page int, filter string) *ListOptions {
	return &ListOptions{PageOptions: &PageOptions{CurrentPage: page}}

}

func applyListOptionsToRequest(opts *ListOptions, req *resty.Request) {
	if opts != nil {
		if opts.PageOptions != nil && opts.CurrentPage > 0 {
			req.SetQueryParam("page", strconv.Itoa(opts.CurrentPage))
		}

		if opts.PageSize > 0 {
			req.SetQueryParam("page_size", strconv.Itoa(opts.PageSize))
		}

	}
}

// listHelper abstracts fetching and pagination for GET endpoints that
// do not require any Ids (top level endpoints).
// When opts (or opts.Page) is nil, all pages will be fetched and
// returned in a single (endpoint-specific)PagedResponse
// opts.results and opts.pages will be updated from the API response
// nolint
func (c *Client) listHelper(ctx context.Context, i interface{}, opts *ListOptions) error {
	var (
		err     error
		pages   int
		results int
		r       *resty.Response
	)

	req := c.R(ctx)
	applyListOptionsToRequest(opts, req)

	switch v := i.(type) {
	case *DomainsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(DomainsPagedResponse{}).Get(v.endpoint(c))); err == nil {
			response, ok := r.Result().(*DomainsPagedResponse)
			if !ok {
				return fmt.Errorf("response is not a *DomainsPagedResponse")
			}
			pages = response.LastPage
			results = response.LastPage
			v.appendData(response)
		}

	default:
		log.Fatalf("listHelper interface{} %+v used", i)
	}

	if err != nil {
		return err
	}

	if opts == nil {
		for page := 2; page <= pages; page++ {
			if err := c.listHelper(ctx, i, &ListOptions{PageOptions: &PageOptions{CurrentPage: page}}); err != nil {
				return err
			}
		}
	} else {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}

		if opts.CurrentPage == 0 {
			for page := 2; page <= pages; page++ {
				opts.CurrentPage = page
				if err := c.listHelper(ctx, i, opts); err != nil {
					return err
				}
			}
		}
		opts.Total = results
		opts.CurrentPage = pages
	}

	return nil
}
