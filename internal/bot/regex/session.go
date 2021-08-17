package regex

import (
	"context"
	"fmt"

	"github.com/erei/avakian/internal/bot"
	"github.com/erei/avakian/internal/pkg/modelsx"
	"github.com/skwair/harmony/discord"
)

type SessionOption func(*Session)

func WithSessionMessage(m *discord.Message) SessionOption {
	return func(s *Session) {
		s.Msg = m
	}
}

func WithSessionClient(c bot.Client) SessionOption {
	return func(s *Session) {
		s.Bot = c
	}
}

func WithSessionGuild(g *modelsx.GuildWithPrefixes) SessionOption {
	return func(s *Session) {
		s.Guild = g
	}
}

type Session struct {
	Msg   *discord.Message
	Bot   bot.Client
	Guild *modelsx.GuildWithPrefixes

	match string
}

func NewSession(opts ...SessionOption) *Session {
	s := new(Session)

	for _, o := range opts {
		o(s)
	}

	return s
}

func (rs *Session) SetMatch(m string) {
	rs.match = m
}

func (rs *Session) Match() string {
	return rs.match
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
