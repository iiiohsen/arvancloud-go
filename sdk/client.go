package sdk

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/go-resty/resty/v2"
)

var envDebug = false

const (
	// APIHost Arvancloud API hostname
	APIHost = "napi.arvancloud.com"
	// APIHostVar environment var to check for alternate API URL
	APIHostVar = "ARVANCLOUD_URL"
	// APIHostCert environment var containing path to CA cert to validate against
	APIHostCert = "ARVANCLOUD_CA"
	// APIVersion Arvancloud API version
	APIVersion = "cdn/4.0"
	// APIVersionVar environment var to check for alternate API Version
	APIVersionVar = "ARVANCLOUD_API_VERSION"
	// APIProto connect to API with http(s)
	APIProto = "https"
	// APIEnvVar environment var to check for API token
	APIEnvVar = "ARVANCLOUD_TOKEN"
	// APISecondsPerPoll how frequently to poll for new Events or Status in WaitFor functions
	APISecondsPerPoll = 3
	// Maximum wait time for retries
	APIRetryMaxWaitTime = time.Duration(30) * time.Second
)

// Client is a wrapper around the Resty client
type Client struct {
	resty             *resty.Client
	userAgent         string
	resources         map[string]*Resource
	debug             bool
	retryConditionals []RetryConditional

	millisecondsPerPoll time.Duration

	baseURL    string
	apiVersion string
	apiProto   string

	Domains       *Resource
	DomainRecords *Resource
	Account       *Resource
}

// R wraps resty's R method
func (c *Client) R(ctx context.Context) *resty.Request {
	return c.resty.R().
		ExpectContentType("application/json").
		SetHeader("Content-Type", "application/json").
		SetContext(ctx).
		SetError(APIError{})
}

// SetAuthHeader sets a custom user-agent for HTTP requests and APIkey
func (c *Client) SetAuthHeader(ua, apkiKey string) *Client {
	c.userAgent = ua
	c.resty.SetHeader("User-Agent", c.userAgent)
	c.resty.SetHeader("Authorization", apkiKey)
	return c
}

func (c *Client) SetUserAgent(apkiKey string) *Client {
	c.resty.SetHeader("User-Agent", apkiKey)
	return c
}

func (c *Client) SetBaseURL(baseURL string) *Client {
	baseURLPath, _ := url.Parse(baseURL)

	c.baseURL = path.Join(baseURLPath.Host, baseURLPath.Path)
	c.apiProto = baseURLPath.Scheme

	c.updateHostURL()

	return c
}

func NewClient(apikey string) (client Client) {

	client.resty = resty.New()

	client.SetAuthHeader(DefaultUserAgent, apikey)
	baseURL, baseURLExists := os.LookupEnv(APIHostVar)

	if baseURLExists {
		client.SetBaseURL(baseURL)
	} else {
		apiVersion, apiVersionExists := os.LookupEnv(APIVersionVar)
		if apiVersionExists {
			client.SetAPIVersion(apiVersion)
		} else {
			client.SetAPIVersion(APIVersion)
		}
	}

	certPath, certPathExists := os.LookupEnv(APIHostCert)

	if certPathExists {
		cert, err := ioutil.ReadFile(certPath)
		if err != nil {
			log.Fatalf("[ERROR] Error when reading cert at %s: %s\n", certPath, err.Error())
		}

		client.SetRootCertificate(certPath)

		if envDebug {
			log.Printf("[DEBUG] Set API root certificate to %s with contents %s\n", certPath, cert)
		}
	}

	client.
		SetRetryWaitTime((1000 * APISecondsPerPoll) * time.Millisecond).
		SetPollDelay(1000 * APISecondsPerPoll).
		SetRetries().
		SetDebug(envDebug)

	addResources(&client)

	return
}

// SetDebug sets the debug on resty's client
func (c *Client) SetDebug(debug bool) *Client {
	c.debug = debug
	c.resty.SetDebug(debug)

	return c
}

// Resource looks up a resource by name
func (c Client) Resource(resourceName string) *Resource {
	selectedResource, ok := c.resources[resourceName]
	if !ok {
		log.Fatalf("Could not find resource named '%s', exiting.", resourceName)
	}

	return selectedResource
}

// nolint
func addResources(client *Client) {
	resources := map[string]*Resource{
		accountName: NewResource(client, accountName, accountEndpoint, false, Account{}, nil),
		domainsName: NewResource(client, domainsName, domainsEndpoint, false, Domain{}, DomainsPagedResponse{}),
	}

	client.resources = resources

	client.Account = resources[accountName]
	client.DomainRecords = resources[domainRecordsName]
	client.Domains = resources[domainsName]

}

func (c *Client) SetRetries() *Client {
	c.
		addRetryConditional(ArvancloudBusyRetryCondition).
		addRetryConditional(tooManyRequestsRetryCondition).
		addRetryConditional(serviceUnavailableRetryCondition).
		addRetryConditional(requestTimeoutRetryCondition).
		SetRetryMaxWaitTime(APIRetryMaxWaitTime)
	configureRetries(c)
	return c
}

func (c *Client) addRetryConditional(retryConditional RetryConditional) *Client {
	c.retryConditionals = append(c.retryConditionals, retryConditional)
	return c
}

// SetRetryMaxWaitTime sets the maximum delay before retrying a request.
func (c *Client) SetRetryMaxWaitTime(max time.Duration) *Client {
	c.resty.SetRetryMaxWaitTime(max)
	return c
}

// SetPollDelay sets the number of milliseconds to wait between events or status polls.
// Affects all WaitFor* functions and retries.
func (c *Client) SetPollDelay(delay time.Duration) *Client {
	c.millisecondsPerPoll = delay
	return c
}

// SetRetryWaitTime sets the default (minimum) delay before retrying a request.
func (c *Client) SetRetryWaitTime(min time.Duration) *Client {
	c.resty.SetRetryWaitTime(min)
	return c
}

// SetRootCertificate adds a root certificate to the underlying TLS client config
func (c *Client) SetRootCertificate(path string) *Client {
	c.resty.SetRootCertificate(path)
	return c
}

// SetAPIVersion sets the version of the API to interface with
func (c *Client) SetAPIVersion(apiVersion string) *Client {
	c.apiVersion = apiVersion

	c.updateHostURL()

	return c
}

func (c *Client) updateHostURL() {
	apiProto := APIProto
	baseURL := APIHost
	apiVersion := APIVersion

	if c.baseURL != "" {
		baseURL = c.baseURL
	}

	if c.apiVersion != "" {
		apiVersion = c.apiVersion
	}

	if c.apiProto != "" {
		apiProto = c.apiProto
	}

	c.resty.SetBaseURL(fmt.Sprintf("%s://%s/%s", apiProto, baseURL, apiVersion))
}
