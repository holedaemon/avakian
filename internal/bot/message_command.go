package bot

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/erei/avakian/internal/database/models"
	"github.com/skwair/harmony/discord"
)

type MessageSession struct {
	Msg *discord.Message
	// Entire message
	Argv []string
	// Without command trigger
	Args   []string
	Prefix string

	Bot *Bot
	Tx  *sql.Tx
}

func (ms *MessageSession) QueryGuild(ctx context.Context) (*models.Guild, error) {
	return ms.Bot.QueryGuild(ctx, ms.Msg.GuildID)
}

func (ms *MessageSession) Reply(ctx context.Context, msg string) error {
	channel := ms.Bot.Client.Channel(ms.Msg.ChannelID)
	_, err := channel.SendMessage(ctx, msg)
	return err
}

func (ms *MessageSession) Replyf(ctx context.Context, msg string, args ...interface{}) error {
	msg = fmt.Sprintf(msg, args...)
	return ms.Reply(ctx, msg)
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
