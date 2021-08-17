package scrape

import (
	"context"
	"net/http"
)

func GetTikTokVideoURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgent)
	return "", nil
}
