package modelsx

import (
	"context"

	"github.com/erei/avakian/internal/database/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type GuildWithPrefix struct {
	models.Guild  `boil:",bind"`
	models.Prefix `boil:",bind"`
}

func GetGuildWithPrefixes(ctx context.Context, exec boil.ContextExecutor, guildID string) (*GuildWithPrefix, error) {
	var gp *GuildWithPrefix

	err := models.NewQuery(
		qm.Select("prefixes.prefix"),
		qm.From("guilds"),
		qm.InnerJoin("prefixes ON prefixes.guild_snowflake = ?", guildID),
	).Bind(ctx, exec, gp)

	if err != nil {
		return nil, err
	}

	return gp, nil
}
