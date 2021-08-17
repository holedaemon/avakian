package message

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/holedaemon/avakian/internal/bot"
	"github.com/holedaemon/avakian/internal/pkg/modelsx"
	"github.com/skwair/harmony/discord"
)

type SessionOption func(s *Session)

func WithSessionClient(b bot.Client) SessionOption {
	return func(s *Session) {
		s.Bot = b
	}
}

func WithSessionTx(tx *sql.Tx) SessionOption {
	return func(s *Session) {
		s.Tx = tx
	}
}

func WithSessionGuild(g *modelsx.GuildWithPrefixes) SessionOption {
	return func(s *Session) {
		s.Guild = g
	}
}

func WithSessionMessage(m *discord.Message) SessionOption {
	return func(s *Session) {
		s.Msg = m
	}
}

type Session struct {
	Msg *discord.Message

	// Entire message
	Argv []string

	// Without command trigger
	Args []string

	// Prefix that triggered the command
	Prefix string

	// DB row for the origin guild
	Guild *modelsx.GuildWithPrefixes

	Bot bot.Client
	Tx  *sql.Tx
}

func NewMessageSession(opts ...SessionOption) *Session {
	ms := new(Session)

	for _, o := range opts {
		o(ms)
	}

	ms.setup()

	return ms
}

func (ms *Session) setup() {
	if ms.Msg == nil {
		panic("msg was never given to MessageSession")
	}

	ms.Prefix = ms.Msg.Content[:1]
	ms.Argv = strings.Split(ms.Msg.Content, " ")
	ms.Args = ms.Argv[1:]
}

func (ms *Session) Reply(ctx context.Context, msg string) error {
	return ms.Bot.SendMessage(ctx, ms.Msg.ChannelID, msg)
}

func (ms *Session) Replyf(ctx context.Context, msg string, args ...interface{}) error {
	msg = fmt.Sprintf(msg, args...)
	return ms.Reply(ctx, msg)
}

func (ms *Session) Client() bot.Client {
	return ms.Bot
}
