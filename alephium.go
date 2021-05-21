package alephium

import (
	"github.com/dghubble/sling"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type AlephiumClient struct {
	oldEndpoint string
	oldClient   *http.Client
	slingClient *sling.Sling
	log         *logrus.Logger
	sleepTime   time.Duration
}

func New(alephiumEndpoint string, log *logrus.Logger) (*AlephiumClient, error) {

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	slingClient := sling.New().Client(client).Base(alephiumEndpoint)
	alephiumClient := &AlephiumClient{
		oldEndpoint: alephiumEndpoint,
		oldClient:   client,
		slingClient: slingClient,
		log:         log,
		sleepTime: 5 * time.Second,
	}

	return alephiumClient, nil
}

func relevantError(e1 error, e2 ErrorDetail) error {
	if e1 != nil {
		return e1
	} else if e2.Detail != "" {
		return e2
	}
	return nil
}
