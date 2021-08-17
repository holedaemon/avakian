package avakian

import (
	"context"
	"strings"

	"github.com/holedaemon/avakian/internal/bot/message"
	"github.com/holedaemon/avakian/internal/database/models"
	"github.com/skwair/harmony/discord"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	cmdGuildctl = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionManageGuild),
		message.WithCommandUsage(guildctlCommands.BuildUsage("guildctl")),
		message.WithCommandFn(cmdGuildctlFn),
	)

	cmdGuildctlOn = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionManageGuild),
		message.WithCommandFn(cmdGuildctlOnFn),
	)

	cmdGuildctlOff = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionManageGuild),
		message.WithCommandFn(cmdGuildctlOffFn),
	)
)

var guildctlCommands = message.NewCommandMap(
	message.WithMapScope(true),
	message.WithMapCommand("on", cmdGuildctlOn),
	message.WithMapCommand("off", cmdGuildctlOff),
)

func cmdGuildctlFn(ctx context.Context, s *message.Session) error {
	return guildctlCommands.ExecuteCommand(ctx, s)
}

var settingsMap = map[string]string{
	"twitterembeds": "embed_twitter_videos",
	"twittervideos": "embed_twitter_videos",
	"te":            "embed_twitter_videos",
	"tv":            "embed_twitter_videos",
}

func toggleSetting(ctx context.Context, s *message.Session, setting string, val bool) (bool, error) {
	guild, err := models.Guilds(qm.Where("guild_snowflake = ?", s.Msg.GuildID)).One(ctx, s.Tx)
	if err != nil {
		return false, err
	}

	switch setting {
	case "embed_twitter_videos":
		guild.EmbedTwitterVideos = val
	}

	if err := guild.Update(ctx, s.Tx, boil.Infer()); err != nil {
		return false, err
	}

	return true, nil
}

func cmdGuildctlOnFn(ctx context.Context, s *message.Session) error {
	arg := strings.ToLower(s.Args[0])

	proc := settingsMap[arg]
	if proc == "" {
		return s.Replyf(ctx, "Unknown setting \"%s\"", arg)
	}

	val := true
	updated, err := toggleSetting(ctx, s, proc, val)
	if err != nil {
		return err
	}

	if !updated {
		//this shouldn't happen
		return s.Reply(ctx, "Unable to update guild setting")
	}

	return s.Replyf(ctx, "Guild setting \"%s\" has been turned on", proc)
}

func cmdGuildctlOffFn(ctx context.Context, s *message.Session) error {
	arg := s.Args[0]

	proc := settingsMap[arg]
	if proc == "" {
		return s.Replyf(ctx, "Unknown setting \"%s\"", arg)
	}

	val := false
	updated, err := toggleSetting(ctx, s, proc, val)
	if err != nil {
		return err
	}

	if !updated {
		return s.Reply(ctx, "Unable to update guild setting")
	}

	return s.Replyf(ctx, "Guild setting \"%s\" has been disabled", proc)
}
