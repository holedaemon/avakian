package reaction

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/holedaemon/avakian/internal/bot"
	"github.com/holedaemon/avakian/internal/pkg/modelsx"
	"github.com/skwair/harmony"
)

type Session struct {
	Reaction *harmony.MessageReaction

	Guild *modelsx.GuildWithPrefixes

	Bot bot.Client
	Tx  *sql.Tx
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
