package tabscanner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	processURL = "https://api.tabscanner.com/%v/process"
	resultURL  = "https://api.tabscanner.com/%v/result/%v"
)

// ClientOption describes a configuration option for the Client.
type ClientOption func(*Client)

// HTTPClientOption allows to inject an HTTP client.
func HTTPClientOption(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// Client describes a TabScanner client.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient initializes a new Client.
func NewClient(apiKey string, options ...ClientOption) *Client {
	c := &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// Process uploads a receipt for processing.
func (c *Client) Process(ctx context.Context, req *ProcessRequest) (*ProcessResponse, error) {
	body, contentType, err := req.toFormBody()
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf(processURL, c.apiKey), body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", contentType)

	httpResp, err := c.httpClient.Do(httpReq.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	resp := &ProcessResponse{}
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Result fetches the processing result for a receipt.
func (c *Client) Result(ctx context.Context, token string) (*ResultResponse, error) {
	httpReq, err := http.NewRequest("GET", fmt.Sprintf(resultURL, c.apiKey, token), nil)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.httpClient.Do(httpReq.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	resp := &ResultResponse{}
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
