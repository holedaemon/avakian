package avakian

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

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

	cmdQuoteGet = message.NewCommand(
		message.WithCommandArgs(1),
		message.WithCommandFn(cmdQuoteGetFn),
		message.WithCommandPermissions(discord.PermissionSendMessages),
	)

	cmdQuoteRemove = message.NewCommand(
		message.WithCommandArgs(1),
		message.WithCommandFn(cmdQuoteRemoveFn),
		message.WithCommandPermissions(discord.PermissionManageMessages),
	)
)

var quoteCommands = message.NewCommandMap(
	message.WithMapScope(true),
	message.WithMapCommand("random", cmdQuoteRandom),
	message.WithMapCommand("get", cmdQuoteGet),
	message.WithMapCommand("remove", cmdQuoteRemove),
	message.WithMapCommand("delete", cmdQuoteRemove),
)

func sendQuote(ctx context.Context, s *message.Session, quote *models.Quote) error {
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

	return sendQuote(ctx, s, quote)
}

func cmdQuoteGetFn(ctx context.Context, s *message.Session) error {
	strIdx := s.Args[0]
	idx, err := strconv.ParseInt(strIdx, 10, 8)
	if err != nil {
		return err
	}

	quote, err := models.Quotes(
		qm.Where("guild_snowflake = ?", s.Msg.GuildID),
		qm.Where("idx = ?", idx),
	).One(ctx, s.Tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.Reply(ctx, "No quote at that index exists")
		}

		return err
	}

	return sendQuote(ctx, s, quote)
}

func cmdQuoteRemoveFn(ctx context.Context, s *message.Session) error {
	strIdx := s.Args[0]
	idx, err := strconv.ParseInt(strIdx, 10, 8)
	if err != nil {
		return err
	}

	err = models.Quotes(
		qm.Where("idx = ?", idx),
		qm.Where("guild_snowflake = ?", s.Msg.GuildID),
	).DeleteAll(ctx, s.Tx)

	if err != nil {
		return err
	}

	return s.Replyf(ctx, "Quote #%d has been removed from the database", idx)
}
