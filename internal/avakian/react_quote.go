package avakian

import (
	"context"

	"github.com/holedaemon/avakian/internal/bot/reaction"
	"github.com/holedaemon/avakian/internal/database/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	reactQuote = reaction.NewCommand(
		reaction.WithCommandFn(reactQuoteFn),
	)
)

// TODO: check if guild has quotes enabled

func reactQuoteFn(ctx context.Context, s *reaction.Session) error {
	exists, err := models.Quotes(qm.Where("message_snowflake = ?", s.Reaction.MessageID)).Exists(ctx, s.Tx)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	var row struct {
		MaxNum null.Int
	}

	models.Quotes(
		qm.Select("max("+models.QuoteColumns.Idx+") AS max_num"),
	).Bind(ctx, s.Tx, &row)

	b := getBot(s)
	msg, err := b.FetchMessage(ctx, s.Reaction.ChannelID, s.Reaction.MessageID)
	if err != nil {
		return err
	}

	q := &models.Quote{
		AuthorSnowflake:  msg.Author.ID,
		QuoterSnowflake:  s.Reaction.UserID,
		Idx:              row.MaxNum.Int + 1,
		MessageSnowflake: s.Reaction.MessageID,
		GuildSnowflake:   s.Reaction.GuildID,
	}

	if err := q.Insert(ctx, s.Tx, boil.Whitelist()); err != nil {
		return err
	}

	return s.Replyf(ctx, "quote added xD")
}
