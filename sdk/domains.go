package sdk

import (
	"context"
	"time"
)

// Domain represents a Domain object
type Domain struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Services struct {
		DNS           string `json:"dns"`
		Cdn           bool   `json:"cdn"`
		CloudSecurity bool   `json:"cloud_security"`
	} `json:"services"`
	DNSCloud           bool      `json:"dns_cloud"`
	PlanLevel          int       `json:"plan_level"`
	PlanDowngradableAt time.Time `json:"plan_downgradable_at"`
	Features           struct {
		EditableVetitiRules    bool `json:"editable_vetiti_rules"`
		EditableDdosRules      bool `json:"editable_ddos_rules"`
		EditableRateLimitRules bool `json:"editable_rate_limit_rules"`
		PackagesForWaf         bool `json:"packages_for_waf"`
		UseNewPlans            bool `json:"use_new_plans"`
		UseHealthCheck         bool `json:"use_health_check"`
		FirewallRuleExpr       bool `json:"firewall_rule_expr"`
		UseNewSslModule        bool `json:"use_new_ssl_module"`
		UseNewLoadBalancer     bool `json:"use_new_load_balancer"`
	} `json:"features"`
	NsKeys             []string  `json:"ns_keys"`
	SmartRoutingStatus string    `json:"smart_routing_status"`
	CurrentNs          []string  `json:"current_ns"`
	Status             string    `json:"status"`
	ParentDomain       bool      `json:"parent_domain"`
	IsPaused           bool      `json:"is_paused"`
	IsSuspended        bool      `json:"is_suspended"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
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
