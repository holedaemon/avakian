package avakian

import (
	"context"

	"github.com/holedaemon/avakian/internal/bot"
	"github.com/holedaemon/avakian/internal/bot/message"
	"github.com/holedaemon/avakian/internal/pkg/snowflake"
	"github.com/skwair/harmony/discord"
)

var cmdSnowflake = message.NewCommand(
	message.WithCommandArgs(1),
	message.WithCommandFn(cmdSnowflakeFn),
	message.WithCommandPermissions(discord.PermissionSendMessages),
)

func cmdSnowflakeFn(ctx context.Context, s *message.Session) error {
	sf := s.Args[0]

	if !snowflake.Valid(sf) {
		return bot.ErrUsage
	}

	t, err := snowflake.Time(sf)
	if err != nil {
		return err
	}

	return s.Replyf(ctx, "The given snowflake was created on <t:%d>", t.Unix())
}
