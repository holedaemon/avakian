package bot

import (
	"context"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/erei/avakian/internal/pkg/snowflake"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
	"mvdan.cc/xurls/v2"
)

var (
	rxu = xurls.Strict()

	regTwitter = &RegexCommand{
		fn: regTwitterFn,
	}
)

func regTwitterFn(ctx context.Context, s *RegexSession) error {
	ur := rxu.FindString(s.Msg.Content)

	if ur == "" {
		return nil
	}

	u, err := url.Parse(ur)
	if err != nil {
		return err
	}

	if u.Host != "twitter.com" {
		return nil
	}

	path := u.Path[1:]
	splits := strings.Split(path, "/")

	if len(splits) < 3 {
		ctxlog.Debug(ctx, "splits less than 3", zap.Any("splits", splits))
		return nil
	}

	id := splits[2]
	i, err := snowflake.AsInt64(id)
	if err != nil {
		return nil
	}

	tweet, res, err := s.Bot.Twitter.Statuses.Show(i, &twitter.StatusShowParams{
		IncludeEntities: twitter.Bool(true),
	})
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		ctxlog.Debug(ctx, "non-ok status", zap.Any("status", res.StatusCode))
	}

	if len(tweet.Entities.Media) < 1 {
		ctxlog.Debug(ctx, "tweet contains no entities")
		return nil
	}

	for _, m := range tweet.ExtendedEntities.Media {
		if len(m.VideoInfo.Variants) > 0 {
			variants := m.VideoInfo.Variants

			sort.Slice(variants, func(i, j int) bool {
				return variants[i].Bitrate > variants[j].Bitrate
			})

			for _, v := range variants {
				if v.ContentType != "video/mp4" {
					continue
				}

				return s.Reply(ctx, v.URL)
			}
		}
	}

	return nil
}
