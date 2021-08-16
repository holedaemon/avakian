package avakian

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/erei/avakian/internal/bot/message"
	"github.com/erei/avakian/internal/database/models"
	"github.com/erei/avakian/internal/pkg/zapx"
	"github.com/skwair/harmony/discord"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

const guildMaxPrefix = 5

var (
	cmdPrefix = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionSendMessages),
		message.WithCommandUsage(prefixCommands.BuildUsage("prefix")),
		message.WithCommandFn(cmdPrefixFn),
	)

	cmdPrefixDelete = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionSendMessages),
		message.WithCommandFn(cmdPrefixRemoveFn),
	)

	cmdPrefixAdd = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionManageGuild),
		message.WithCommandFn(cmdPrefixAddFn),
	)

	cmdPrefixList = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionSendMessages),
		message.WithCommandFn(cmdPrefixListFn),
	)

	prefixCommands = message.NewCommandMap(
		message.WithMapScope(true),
		message.WithMapCommand("add", cmdPrefixAdd),
		message.WithMapCommand("delete", cmdPrefixDelete),
		message.WithMapCommand("remove", cmdPrefixDelete),
		message.WithMapCommand("list", cmdPrefixList),
	)
)

func cmdPrefixFn(ctx context.Context, s *message.Session) error {
	return prefixCommands.ExecuteCommand(ctx, s)
}

func cmdPrefixAddFn(ctx context.Context, s *message.Session) error {
	b := getBot(s)

	newPrefix := s.Args[0]
	if len(newPrefix) > 1 {
		return s.Reply(ctx, "For the sake of simplicity, prefixes can only be a single character in length")
	}

	if newPrefix == b.DefaultPrefix {
		return s.Reply(ctx, "You can't act upon the global prefix, doofus")
	}

	exists, err := models.Prefixes(qm.Where("guild_snowflake = ? AND prefix = LOWER(?)", s.Msg.GuildID, newPrefix)).Exists(ctx, s.Tx)
	if err != nil {
		return err
	}

	if exists {
		return s.Reply(ctx, "That's literally already a prefix here")
	}

	count, err := models.Prefixes(qm.Where("guild_snowflake = ?", s.Msg.GuildID)).Count(ctx, s.Tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			count = 0
		} else {
			return err
		}
	}

	if count >= guildMaxPrefix {
		return s.Replyf(ctx, "You ain't allowed more than %d prefixes and that's final", guildMaxPrefix)
	}

	pref := &models.Prefix{
		GuildSnowflake: s.Msg.GuildID,
		Prefix:         newPrefix,
	}

	if err := pref.Insert(ctx, s.Tx, boil.Infer()); err != nil {
		return err
	}

	return s.Replyf(ctx, "Character %s will henceforth be accepted as a prefix", newPrefix)
}

func cmdPrefixRemoveFn(ctx context.Context, s *message.Session) error {
	b := getBot(s)

	oldPrefix := s.Args[0]
	if len(oldPrefix) > 1 {
		return s.Reply(ctx, "Prefixes aren't more than a single character in length")
	}

	if oldPrefix == b.DefaultPrefix {
		return s.Reply(ctx, "You can't act upon the global prefix, doofus")
	}

	pref, err := models.Prefixes(qm.Where("guild_snowflake = ? AND prefix = LOWER(?)", s.Msg.GuildID, oldPrefix)).One(ctx, s.Tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctxlog.Debug(ctx, "prefix doesn't exist")
			return s.Reply(ctx, "That isn't a prefix for this guild, dude")
		}

		ctxlog.Error(ctx, "error removing prefix", zap.Error(err), zapx.Guild(s.Msg.GuildID))
		return err
	}

	if err := pref.Delete(ctx, s.Tx); err != nil {
		return err
	}

	ctxlog.Info(ctx, "removed prefix in guild", zap.String("prefix", oldPrefix), zapx.Guild(s.Msg.GuildID))
	return s.Reply(ctx, "Prefix has been removed")
}

func cmdPrefixListFn(ctx context.Context, s *message.Session) error {
	var sb strings.Builder
	sb.WriteString("Registered prefixes:\n```")

	prefixes, err := models.Prefixes(qm.Where("guild_snowflake = ?", s.Msg.GuildID)).All(ctx, s.Tx)
	if err != nil {
		return err
	}

	if len(prefixes) == 0 {
		return s.Reply(ctx, "Guild hasn't added any prefixes")
	}

	for _, p := range prefixes {
		sb.WriteString(p.Prefix + "\n")
	}

	sb.WriteString("```")
	return s.Reply(ctx, sb.String())
}
