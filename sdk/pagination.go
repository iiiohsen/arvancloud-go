package sdk

/**
 * Pagination and Filtering types and helpers
 */

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// PageOptions are the pagination parameters for List endpoints
type PageOptions struct {
	Page    int `url:"page,omitempty" json:"page"`
	Pages   int `url:"pages,omitempty" json:"pages"`
	Results int `url:"results,omitempty" json:"results"`
}

// ListOptions are the pagination and filtering (TODO) parameters for endpoints
type ListOptions struct {
	*PageOptions
	PageSize int
	Filter   string
}

// NewListOptions simplified construction of ListOptions using only
// the two writable properties, Page and Filter
func NewListOptions(page int, filter string) *ListOptions {
	return &ListOptions{PageOptions: &PageOptions{Page: page}, Filter: filter}
}

func applyListOptionsToRequest(opts *ListOptions, req *resty.Request) {
	if opts != nil {
		if opts.PageOptions != nil && opts.Page > 0 {
			req.SetQueryParam("page", strconv.Itoa(opts.Page))
		}

		if opts.PageSize > 0 {
			req.SetQueryParam("page_size", strconv.Itoa(opts.PageSize))
		}

		if len(opts.Filter) > 0 {
			req.SetHeader("X-Filter", opts.Filter)
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
			pages = response.Pages
			results = response.Results
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
			if err := c.listHelper(ctx, i, &ListOptions{PageOptions: &PageOptions{Page: page}}); err != nil {
				return err
			}
		}
	} else {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}

		if opts.Page == 0 {
			for page := 2; page <= pages; page++ {
				opts.Page = page
				if err := c.listHelper(ctx, i, opts); err != nil {
					return err
				}
			}
		}
		opts.Results = results
		opts.Pages = pages
	}

	return nil
}

// listHelperWithID abstracts fetching and pagination for GET endpoints that
// require an Id (second level endpoints).
// When opts (or opts.Page) is nil, all pages will be fetched and
// returned in a single (endpoint-specific)PagedResponse
// opts.results and opts.pages will be updated from the API response
// nolint
func (c *Client) listHelperWithID(ctx context.Context, i interface{}, idRaw interface{}, opts *ListOptions) error {
	var (
		err     error
		pages   int
		results int
		r       *resty.Response
	)

	req := c.R(ctx)
	applyListOptionsToRequest(opts, req)

	id, _ := idRaw.(int)

	switch v := i.(type) {
	case *DomainRecordsPagedResponse:
		if r, err = coupleAPIErrors(req.SetResult(DomainRecordsPagedResponse{}).Get(v.endpointWithID(c, id))); err == nil {
			response, ok := r.Result().(*DomainRecordsPagedResponse)
			if !ok {
				return fmt.Errorf("response is not a *DomainRecordsPagedResponse")
			}
			pages = response.Pages
			results = response.Results
			v.appendData(response)
		}

	default:
		log.Fatalf("Unknown listHelperWithID interface{} %T used", i)
	}

	if err != nil {
		return err
	}

	if opts == nil {
		for page := 2; page <= pages; page++ {
			if err := c.listHelperWithID(ctx, i, id, &ListOptions{PageOptions: &PageOptions{Page: page}}); err != nil {
				return err
			}
		}
	} else {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}
		if opts.Page == 0 {
			for page := 2; page <= pages; page++ {
				opts.Page = page
				if err := c.listHelperWithID(ctx, i, id, opts); err != nil {
					return err
				}
			}
		}
		opts.Results = results
		opts.Pages = pages
	}

	return nil
}
