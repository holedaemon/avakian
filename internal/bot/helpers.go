package bot

import (
	"context"
	"strings"

	"github.com/erei/avakian/internal/pkg/zapx"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func stringInSlice(want string, sl []string) bool {
	for _, s := range sl {
		if strings.EqualFold(want, s) {
			return true
		}
	}

	return false
}

func buildUsage(prefix, command string, commands interface{}) string {
	switch commands := commands.(type) {
	case map[string]*MessageCommand:
		var sb strings.Builder

		sb.WriteString(prefix + command + " <")

		i := 1
		for k := range commands {
			if i == len(commands) {
				sb.WriteString(k + ">")
				break
			}

			sb.WriteString(k + "|")

			i++
		}

		return sb.String()
	default:
		panic("invalid type passed to buildUsage()")
	}
}

func (b *Bot) runMessageSubcommand(ctx context.Context, s *MessageSession, subCommands map[string]*MessageCommand, usage func() error) error {
	subSess := *s
	subSess.Args = s.Args[1:]
	sub := s.Args[0]

	cmd := subCommands[sub]
	if cmd == nil {
		return usage()
	}

	if !stringInSlice(s.Msg.Author.ID, b.Admins) {
		p, err := s.Bot.FetchMemberPermissions(ctx, s.Msg.GuildID, s.Msg.ChannelID, s.Msg.Author.ID)
		if err != nil {
			return err
		}

		if !cmd.HasPermission(p) {
			ctxlog.Debug(ctx, "member lacks permission to run command", zapx.Member(s.Msg.Author.ID))
			return nil
		}
	} else {
		ctxlog.Debug(ctx, "user is an admin", zap.String("user_id", s.Msg.Author.ID))
	}

	if err := cmd.Execute(ctx, &subSess); err != nil {
		return err
	}

	return nil
}
