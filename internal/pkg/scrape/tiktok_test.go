package scrape

import (
	"context"
	"net/http"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestCanScrapeTikTok(t *testing.T) {
	c := New(WithHTTPClient(&http.Client{Timeout: time.Second * 10}))

	ctx := context.Background()
	video, err := c.GetTikTokVideoURL(ctx, "https://vm.tiktok.com/ZMdT1pYBg/")
	assert.NilError(t, err, "getting url")

	assert.Assert(t, video != "", "video url ")
}
