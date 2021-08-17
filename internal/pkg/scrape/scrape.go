package scrape

import (
	"context"
	"errors"
	"net/http"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0"

var (
	ErrHTTPStatus = errors.New("scrape: non-OK http status code")
)

type Client struct {
	cli *http.Client
}

func New(opts ...Option) *Client {
	c := new(Client)

	for _, o := range opts {
		o(c)
	}

	return c
}

func (c *Client) makeBaseRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	return req, nil
}
