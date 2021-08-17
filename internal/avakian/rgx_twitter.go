package avakian

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/erei/avakian/internal/bot/regex"
	"github.com/erei/avakian/internal/pkg/httpx"
	"github.com/erei/avakian/internal/pkg/snowflake"
	"github.com/zikaeroh/ctxlog"
)

var regTwitter = regex.NewCommand(
	regex.WithCommandFn(regTwitterFn),
)

var (
	reTwitterSf = regexp.MustCompile(`\d{1,20}`)
)

func regTwitterFn(ctx context.Context, s *regex.Session) error {
	if !s.Guild.EmbedTwitterVideos {
		return fmt.Errorf("%w: guild has embeds disabled", regex.ErrCondition)
	}

	link := s.Match()

	sf := reTwitterSf.FindString(link)
	if sf == "" {
		return fmt.Errorf("%w: could not find snowflake in message", regex.ErrCondition)
	}

	sfi, err := snowflake.AsInt64(sf)
	if err != nil {
		return err
	}

	b := getBot(s)
	tweet, res, err := b.Twitter.Statuses.Show(sfi, &twitter.StatusShowParams{
		IncludeEntities: twitter.Bool(true),
		TweetMode:       "extended",
	})
	if err != nil {
		return err
	}

	if res.StatusCode < http.StatusOK && res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("%w: %d", httpx.ErrStatusCode, res.StatusCode)
	}

	if len(tweet.ExtendedEntities.Media) < 1 {
		ctxlog.Debug(ctx, "tweet doesn't contain video")
		return fmt.Errorf("%w: tweet does not contain a video", regex.ErrCondition)
	}

	for _, m := range tweet.ExtendedEntities.Media {
		if m.Type != "video" {
			continue
		}

		variants := m.VideoInfo.Variants
		sort.Slice(variants, func(i, j int) bool {
			return variants[i].Bitrate > variants[j].Bitrate
		})

		return s.Reply(ctx, variants[0].URL)
	}

	return nil
}
