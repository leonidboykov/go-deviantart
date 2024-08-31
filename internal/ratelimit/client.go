package ratelimit

import (
	"errors"
	"net/http"
	"time"
)

const (
	defaultMaxRetries     = 100
	defaultBackoffTimeout = 500 * time.Millisecond
)

var errMaxRetries = errors.New("max retries exceeded")

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClient struct {
	client Doer
}

func NewHTTPClient(client Doer) *HTTPClient {
	return &HTTPClient{client: client}
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	backoffTimeout := defaultBackoffTimeout
	for i := 0; i < defaultMaxRetries; i++ {
		resp, err := c.client.Do(req)
		if err != nil {
			return resp, err
		}
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, err
		}
		time.Sleep(backoffTimeout)
		backoffTimeout *= 2
	}
	return nil, errMaxRetries
}
