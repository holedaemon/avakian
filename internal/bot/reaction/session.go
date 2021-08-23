package reaction

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/holedaemon/avakian/internal/bot"
	"github.com/holedaemon/avakian/internal/pkg/modelsx"
	"github.com/skwair/harmony"
)

type SessionOption func(*Session)

func WithSessionReaction(r *harmony.MessageReaction) SessionOption {
	return func(s *Session) {
		s.Reaction = r
	}
}

func WithSessionGuild(g *modelsx.GuildWithPrefixes) SessionOption {
	return func(s *Session) {
		s.Guild = g
	}
}

func WithSessionBot(b bot.Client) SessionOption {
	return func(s *Session) {
		s.Bot = b
	}
}

func WithSessionTx(tx *sql.Tx) SessionOption {
	return func(s *Session) {
		s.Tx = tx
	}
}

type Session struct {
	Reaction *harmony.MessageReaction

	Guild *modelsx.GuildWithPrefixes

	Bot bot.Client
	Tx  *sql.Tx
}

func NewSession(opts ...SessionOption) *Session {
	s := new(Session)

	for _, o := range opts {
		o(s)
	}

	return s
}

func (s *Session) Client() bot.Client {
	return s.Bot
}

func (s *Session) Reply(ctx context.Context, msg string) error {
	return s.Bot.SendMessage(ctx, s.Reaction.ChannelID, msg)
}

func (s *Session) Replyf(ctx context.Context, msg string, args ...interface{}) error {
	msg = fmt.Sprintf(msg, args...)
	return s.Reply(ctx, msg)
}

func (s *Session) Mention(ctx context.Context, msg string) error {
	return s.Replyf(ctx, "<@%s>, %s", s.Reaction.UserID, msg)
}
