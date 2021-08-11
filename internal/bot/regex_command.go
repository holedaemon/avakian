package bot

import (
	"context"
	"fmt"

	"github.com/skwair/harmony/discord"
)

type RegexCommand struct {
	fn func(context.Context, *RegexSession) error
}

func (rg *RegexCommand) Execute(ctx context.Context, s Session) error {
	sess, ok := s.(*RegexSession)
	if !ok {
		panic("wrong session passed to RegexCommand")
	}

	return rg.fn(ctx, sess)
}

// no-op
func (rg *RegexCommand) Usage(context.Context, Session) error {
	return nil
}

type RegexSession struct {
	Msg *discord.Message
	Bot *Bot
}

func (rs *RegexSession) Reply(ctx context.Context, msg string) error {
	channel := rs.Bot.Client.Channel(rs.Msg.ChannelID)
	_, err := channel.SendMessage(ctx, msg)
	return err
}

func (rs *RegexSession) Replyf(ctx context.Context, msg string, args ...interface{}) error {
	msg = fmt.Sprintf(msg, args...)
	return rs.Reply(ctx, msg)
}
