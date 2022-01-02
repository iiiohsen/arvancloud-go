package sdk

import "context"

// Domain represents a Domain object
type Domain struct {
	//	This Domain's unique ID
	ID int `json:"id"`

	// The domain this Domain represents. These must be unique in our system; you cannot have two Domains representing the same domain.
	Domain string `json:"domain"`
}

// ListDomains lists Domains
func (c *Client) ListDomains(ctx context.Context, opts *ListOptions) ([]Domain, error) {
	response := DomainsPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// appendData appends Domains when processing paginated Domain responses
func (resp *DomainsPagedResponse) appendData(r *DomainsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// DomainsPagedResponse represents a paginated Domain API response
type DomainsPagedResponse struct {
	*PageOptions
	Data []Domain `json:"data"`
}

// endpoint gets the endpoint URL for Domain
func (DomainsPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.Domains.Endpoint()
	if err != nil {
		panic(err)
	}

	return endpoint
}
