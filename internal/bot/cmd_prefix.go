package bot

import (
	"context"
	"database/sql"
	"errors"
	"strings"

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
	cmdAddPrefix = &MessageCommand{
		permissions: discord.PermissionManageGuild,
		fn:          cmdAddPrefixFn,
	}

	cmdRemovePrefix = &MessageCommand{
		permissions: discord.PermissionManageGuild,
		fn:          cmdRemovePrefixFn,
	}

	cmdListPrefix = &MessageCommand{
		permissions: discord.PermissionSendMessages,
		fn:          cmdListPrefixFn,
	}
)

func cmdAddPrefixFn(ctx context.Context, s *MessageSession) error {
	if len(s.Args) == 0 {
		return s.Reply(ctx, "At least one argument is required")
	}

	newPrefix := s.Args[0]
	if len(newPrefix) > 1 {
		return s.Reply(ctx, "For the sake of simplicity, prefixes can only be a single character in length")
	}

	if newPrefix == s.Bot.DefaultPrefix {
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

func cmdRemovePrefixFn(ctx context.Context, s *MessageSession) error {
	if len(s.Args) == 0 {
		return s.Reply(ctx, "At least 1 argument is required")
	}

	oldPrefix := s.Args[0]
	if len(oldPrefix) > 1 {
		return s.Reply(ctx, "Prefixes aren't more than a single character in length")
	}

	if oldPrefix == s.Bot.DefaultPrefix {
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

func cmdListPrefixFn(ctx context.Context, s *MessageSession) error {
	var sb strings.Builder
	sb.WriteString("Registered prefixes:\n```")

	prefixes, err := models.Prefixes(qm.Where("guild_snowflake = ?", s.Msg.GuildID)).All(ctx, s.Tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctxlog.Debug(ctx, "guild has no prefixes")
			return s.Reply(ctx, "Guild hasn't added any prefixes")
		}

		return err
	}

	for _, p := range prefixes {
		sb.WriteString(p.Prefix + "\n")
	}

	sb.WriteString("```")
	return s.Reply(ctx, sb.String())
}
