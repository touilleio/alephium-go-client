package alephium

import (
	"github.com/dghubble/sling"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Client struct {
	endpointURI string
	apiToken    string
	oldClient   *http.Client
	slingClient *sling.Sling
	log         *logrus.Logger
	sleepTime   time.Duration
}

const (
	ApiKeyHeader = "X-API-KEY"
)

func New(alephiumEndpoint string, log *logrus.Logger) (*Client, error) {
	return NewWithApiKey(alephiumEndpoint, "", log)
}

func NewWithApiKey(alephiumEndpoint string, apiKey string, log *logrus.Logger) (*Client, error) {

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	slingClient := sling.New().Client(client).Base(alephiumEndpoint)

	if apiKey != "" {
		slingClient = slingClient.Add(ApiKeyHeader, apiKey)
	}

	alephiumClient := &Client{
		endpointURI: alephiumEndpoint,
		oldClient:   client,
		slingClient: slingClient,
		log:         log,
		sleepTime:   5 * time.Second,
	}

	return alephiumClient, nil
}

func (a *Client) String() string {
	return a.endpointURI
}

func relevantError(e1 error, e2 ErrorDetail) error {
	if e1 != nil {
		return e1
	} else if e2.Detail != "" {
		return e2
	}
	return nil
}
