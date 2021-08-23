package modelsx

import (
	"context"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/types"
)

type GuildWithPrefixes struct {
	ID                 int64             `boil:"id"`
	GuildSnowflake     string            `boil:"guild_snowflake"`
	EmbedTwitterVideos bool              `boil:"embed_twitter_videos"`
	DoQuotes           bool              `boil:"do_quotes"`
	CreatedAt          time.Time         `boil:"created_at"`
	UpdatedAt          time.Time         `boil:"updated_at"`
	Prefixes           types.StringArray `boil:"prefixes"`
}

func (g *GuildWithPrefixes) HasPrefix(p string) bool {
	for _, pr := range g.Prefixes {
		if pr == p {
			return true
		}
	}

	return false
}

func GetGuildWithPrefixes(ctx context.Context, exec boil.ContextExecutor, guildID string) (*GuildWithPrefixes, error) {
	gp := new(GuildWithPrefixes)

	err := queries.Raw("SELECT guilds.*, array_agg(prefixes.prefix) FILTER (WHERE prefixes.prefix IS NOT NULL) AS prefixes FROM guilds LEFT JOIN prefixes ON guilds.guild_snowflake = prefixes.guild_snowflake WHERE guilds.guild_snowflake = $1 GROUP BY guilds.id;", guildID).Bind(ctx, exec, gp)
	if err != nil {
		return nil, err
	}

	return gp, nil
}
