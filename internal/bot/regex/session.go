package regex

import (
	"context"
	"fmt"

	"github.com/erei/avakian/internal/bot"
	"github.com/skwair/harmony/discord"
)

type Session struct {
	Msg *discord.Message
	Bot bot.Client
}

func (rs *Session) Reply(ctx context.Context, msg string) error {
	return rs.Bot.SendMessage(ctx, rs.Msg.ChannelID, msg)
}

func (rs *Session) Replyf(ctx context.Context, msg string, args ...interface{}) error {
	msg = fmt.Sprintf(msg, args...)
	return rs.Reply(ctx, msg)
}

func (rs *Session) Client() bot.Client {
	return rs.Bot
}
