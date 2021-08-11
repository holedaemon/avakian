package bot

import (
	"context"
	"strings"

	"github.com/erei/avakian/internal/database/models"
	"github.com/skwair/harmony/discord"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	cmdGuildctl = &MessageCommand{
		fn:          cmdGuildctlFn,
		permissions: discord.PermissionManageGuild,
		usage:       buildUsage("guildctl", guildctlCommands),
	}

	cmdGuildctlOn = &MessageCommand{
		fn:          cmdGuildctlOnFn,
		permissions: discord.PermissionManageGuild,
	}

	cmdGuildctlOff = &MessageCommand{
		fn:          cmdGuildctlOffFn,
		permissions: discord.PermissionManageGuild,
	}
)

var guildctlCommands = messageCommandMap{
	"on":  cmdGuildctlOn,
	"off": cmdGuildctlOff,
}

func cmdGuildctlFn(ctx context.Context, s *MessageSession) error {
	if len(s.Args) == 0 {
		return ErrUsage
	}

	return guildctlCommands.ExecuteSubCommand(ctx, s)
}

var settingsMap = map[string]string{
	"twitterembeds": "embed_twitter_videos",
	"twittervideos": "embed_twitter_videos",
	"te":            "embed_twitter_videos",
	"tv":            "embed_twitter_videos",
}

func toggleSetting(ctx context.Context, s *MessageSession, setting string, val bool) (bool, error) {
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

func cmdGuildctlOnFn(ctx context.Context, s *MessageSession) error {
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

func cmdGuildctlOffFn(ctx context.Context, s *MessageSession) error {
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
