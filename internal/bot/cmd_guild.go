package bot

import (
	"context"
	"strings"

	"github.com/skwair/harmony/discord"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

var (
	cmdGuildctl = &MessageCommand{
		fn:          cmdGuildctlFn,
		permissions: discord.PermissionManageGuild,
	}

	cmdGuildctlOn = &MessageCommand{
		fn:          cmdGuildctlOnFn,
		permissions: discord.PermissionManageGuild,
	}

	cmdGuildctlOff = &MessageCommand{
		fn:          cmdGuildctlOffFn,
		permissions: discord.PermissionManageGuild,
	}

	cmdGuildctlToggle = &MessageCommand{
		fn:          cmdGuildctlToggleFn,
		permissions: discord.PermissionManageGuild,
	}
)

var guildctlCommands = map[string]*MessageCommand{
	"on":     cmdGuildctlOn,
	"off":    cmdGuildctlOff,
	"toggle": cmdGuildctlToggle,
}

func cmdGuildctlFn(ctx context.Context, s *MessageSession) error {
	usage := func() error {
		return s.Replyf(ctx, "Usage: `%s`", buildUsage(s.Prefix, "guildctl", guildctlCommands))
	}

	if len(s.Args) == 0 {
		return usage()
	}

	subSess := *s
	subSess.Args = s.Args[1:]
	sub := s.Args[0]

	cmd := guildctlCommands[sub]
	if cmd == nil {
		return usage()
	}

	p, err := s.Bot.FetchMemberPermissions(ctx, s.Msg.GuildID, s.Msg.ChannelID, s.Msg.Author.ID)
	if err != nil {
		return err
	}

	if !cmd.HasPermission(p) {
		return nil
	}

	if err := cmd.Execute(ctx, &subSess); err != nil {
		return err
	}

	return nil
}

var settingsMap = map[string]string{
	"twitterembeds": "embed_twitter_videos",
	"twittervideos": "embed_twitter_videos",
	"te":            "embed_twitter_videos",
	"tv":            "embed_twitter_videos",
}

func toggleSetting(ctx context.Context, s *MessageSession, setting string, val *bool) (bool, bool, error) {
	if settingsMap[setting] == "" {
		return false, false, nil
	}

	if val == nil {
		var row struct {
			Setting bool
		}

		query := queries.Raw(
			"UPDATE guilds SET $1 = !$1 AND updated_at = NOW() WHERE guild_snowflake = $2 RETURNING $1 AS setting;",
			setting,
			s.Msg.GuildID,
		)

		if err := query.Bind(ctx, s.Tx, &row); err != nil {
			return false, false, err
		}

		return true, row.Setting, nil
	}

	query := queries.Raw(
		"UPDATE guilds SET $1 = $2 AND updated_at = NOW() WHERE guild_snowflake = $3;",
		setting,
		*val,
		s.Msg.GuildID,
	)

	_, err := query.ExecContext(ctx, s.Tx)

	if err != nil {
		return false, false, err
	}

	return true, false, nil
}

func cmdGuildctlOnFn(ctx context.Context, s *MessageSession) error {
	arg := strings.ToLower(s.Args[0])

	proc := settingsMap[arg]
	if proc == "" {
		return s.Replyf(ctx, "Unknown setting \"%s\"", arg)
	}

	var val *bool
	*val = true
	updated, _, err := toggleSetting(ctx, s, proc, val)
	if err != nil {
		return err
	}

	if !updated {
		//this shouldn't happen
		return s.Reply(ctx, "Unable to update guild setting")
	}

	return s.Replyf(ctx, "Guild setting \"%s\" has been turned on", arg)
}

func cmdGuildctlOffFn(ctx context.Context, s *MessageSession) error {
	arg := s.Args[0]

	proc := settingsMap[arg]
	if proc == "" {
		return s.Replyf(ctx, "Unknown setting \"%s\"", arg)
	}

	var val *bool
	*val = false
	updated, _, err := toggleSetting(ctx, s, proc, val)
	if err != nil {
		return err
	}

	if !updated {
		return s.Reply(ctx, "Unable to update guild setting")
	}

	return s.Replyf(ctx, "Guild setting \"%s\" has been disabled", arg)
}

func cmdGuildctlToggleFn(ctx context.Context, s *MessageSession) error {
	arg := s.Args[0]

	proc := settingsMap[arg]
	if proc == "" {
		return s.Replyf(ctx, "Unknown setting \"%s\"", arg)
	}

	updated, val, err := toggleSetting(ctx, s, proc, nil)
	if err != nil {
		return err
	}

	if !updated {
		return s.Reply(ctx, "Unable to update guild setting")
	}

	return s.Replyf(ctx, "Guild setting \"%s\" has been set to %t", proc, val)
}
