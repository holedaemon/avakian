package bot

import (
	"context"

	"github.com/skwair/harmony/discord"
)

type MessageSession struct {
	Msg  *discord.Message
	Argv []string
	Bot  *Bot
}

func (ms *MessageSession) Reply(ctx context.Context, msg string) error {
	channel := ms.Bot.Client.Channel(ms.Msg.ChannelID)
	_, err := channel.SendMessage(ctx, msg)
	return err
}

type MessageCommand struct {
	permissions int
	fn          func(context.Context, *MessageSession) error
}

// NOTE: check that this is correct
func (mc *MessageCommand) HasPermission(p int) bool {
	if p == discord.PermissionAdministrator {
		return true
	}

	return (p & mc.permissions) != 0
}

func (mc *MessageCommand) Execute(ctx context.Context, s Session) error {
	ms, ok := s.(*MessageSession)
	if !ok {
		panic("wrong session type given")
	}

	return mc.fn(ctx, ms)
}
