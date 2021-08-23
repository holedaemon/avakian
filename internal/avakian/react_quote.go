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

func reactQuoteFn(ctx context.Context, s *reaction.Session) error {
	if !s.Guild.DoQuotes {
		return nil
	}

	b := getBot(s)
	msg, err := b.FetchMessage(ctx, s.Reaction.ChannelID, s.Reaction.MessageID)
	if err != nil {
		return err
	}

	if msg.Author.ID == s.Reaction.UserID {
		return s.Mention(ctx, "self-quoting is disallowed. Self-absorbed much?")
	}

	if msg.Author.ID == b.Client.State.Me().ID {
		return nil
	}

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

	var num int
	if row.MaxNum.Valid {
		num = row.MaxNum.Int + 1
	} else {
		num = 1
	}

	q := &models.Quote{
		AuthorSnowflake:  msg.Author.ID,
		QuoterSnowflake:  s.Reaction.UserID,
		Idx:              num,
		MessageSnowflake: s.Reaction.MessageID,
		GuildSnowflake:   s.Reaction.GuildID,
	}

	if err := q.Insert(ctx, s.Tx, boil.Infer()); err != nil {
		return err
	}

	quoter, err := b.FetchMember(ctx, s.Reaction.UserID, s.Reaction.GuildID)
	if err != nil {
		return err
	}

	author, err := b.FetchMember(ctx, msg.Author.ID, s.Reaction.GuildID)
	if err != nil {
		return err
	}

	return s.Replyf(ctx, "%s quoted a message from %s %s",
		fullUsername(quoter),
		fullUsername(author),
		jumpLinkFromSession(s),
	)
}
