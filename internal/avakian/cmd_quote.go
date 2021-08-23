package avakian

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/holedaemon/avakian/internal/bot/message"
	"github.com/holedaemon/avakian/internal/database/models"
	"github.com/holedaemon/avakian/internal/pkg/snowflake"
	"github.com/skwair/harmony/discord"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	cmdQuote = message.NewCommand(
		message.WithCommandArgs(1),
		message.WithCommandFn(cmdQuoteFn),
		message.WithCommandPermissions(discord.PermissionSendMessages),
	)

	cmdQuoteRandom = message.NewCommand(
		message.WithCommandFn(cmdQuoteRandomFn),
		message.WithCommandPermissions(discord.PermissionSendMessages),
	)
)

var quoteCommands = message.NewCommandMap(
	message.WithMapScope(true),
	message.WithMapCommand("random", cmdQuoteRandom),
)

func cmdQuoteFn(ctx context.Context, s *message.Session) error {
	return quoteCommands.ExecuteCommand(ctx, s)
}

func cmdQuoteRandomFn(ctx context.Context, s *message.Session) error {
	quote, err := models.Quotes(
		qm.OrderBy("random()"),
		qm.Where("guild_snowflake = ?", s.Msg.GuildID),
	).One(ctx, s.Tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.Reply(ctx, "This guild doesn't have any quotes")
		}

		return err
	}

	b := getBot(s)
	msg, err := b.FetchMessage(ctx, quote.ChannelSnowflake, quote.MessageSnowflake)
	if err != nil {
		return err
	}

	return s.Replyf(ctx, "> %s\n― %s, %s\n%s",
		msg.Content,
		fullUsernameFromMessage(msg),
		snowflake.MarkdownTime(msg.ID),
		fmt.Sprintf("<"+jumpLinkURL+">", msg.GuildID, msg.ChannelID, msg.ID),
	)
}
