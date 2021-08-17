package scrape

import "net/http"

type Option func(*Client)

func WithHTTPClient(cli *http.Client) Option {
	return func(c *Client) {
		c.cli = cli
	}
}
