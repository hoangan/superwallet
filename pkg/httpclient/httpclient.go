package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	url    string
}

// NewHttpClient creates a new HttpClient with a given url
func NewHttpClient(url string) *Client {
	return &Client{
		client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   time.Second,
				ResponseHeaderTimeout: time.Second,
			},
		},
		url: url,
	}
}

func (c *Client) Post(body []byte) ([]byte, error) {
	resp, err := c.client.Post(c.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// closing the response body is important to prevent resource leaks
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
