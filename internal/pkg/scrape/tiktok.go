package scrape

import (
	"context"
	"fmt"
	"net/http"

	"github.com/antchfx/htmlquery"
)

func (c *Client) GetTikTokVideoURL(ctx context.Context, url string) (string, error) {
	req, err := c.makeBaseRequest(ctx, url)
	if err != nil {
		return "", err
	}

	res, err := c.cli.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK && res.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("%w: %d", ErrHTTPStatus, res.StatusCode)
	}

	doc, err := htmlquery.Parse(res.Body)
	if err != nil {
		return "", err
	}

	node, err := htmlquery.Query(doc, `//video`)
	if err != nil {
		return "", err
	}

	if node == nil {
		return "", nil
	}

	for _, attr := range node.Attr {
		if attr.Key == "src" {
			return attr.Val, nil
		}
	}

	return "", nil
}
