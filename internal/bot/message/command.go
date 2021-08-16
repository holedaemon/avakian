package message

import (
	"context"

	"github.com/erei/avakian/internal/bot"
	"github.com/skwair/harmony/discord"
)

type CommandOption func(*Command)

func WithCommandPermissions(p int) CommandOption {
	return func(c *Command) {
		c.permissions = p
	}
}

func WithCommandArgs(a int) CommandOption {
	return func(c *Command) {
		c.minArgs = a
	}
}

func WithCommandUsage(us string) CommandOption {
	return func(c *Command) {
		c.usage = us
	}
}

func WithCommandFn(fn func(context.Context, *Session) error) CommandOption {
	return func(c *Command) {
		c.fn = fn
	}
}

type Command struct {
	permissions int
	minArgs     int
	usage       string
	fn          func(context.Context, *Session) error
}

func NewCommand(opts ...CommandOption) *Command {
	c := new(Command)

	for _, o := range opts {
		o(c)
	}

	return c
}

func (mc *Command) HasPermission(p int) bool {
	if p == -1 {
		return true
	}

	if p == discord.PermissionAdministrator {
		return true
	}

	return (p & mc.permissions) != 0
}

func (mc *Command) Execute(ctx context.Context, sess bot.Session) error {
	s, ok := sess.(*Session)
	if !ok {
		panic("wrong session type passed to command")
	}

	p, err := s.Bot.CheckPermission(ctx, sess)
	if err != nil {
		return err
	}

	if !mc.HasPermission(p) {
		return bot.ErrPermission
	}

	return mc.fn(ctx, s)
}

func (mc *Command) Usage(ctx context.Context, s bot.Session) error {
	if mc.usage == "" {
		return s.Reply(ctx, "Example usage was not provided for this command")
	}

	sess, ok := s.(*Session)
	if !ok {
		panic("wrong session type passed to command")
	}

	return s.Replyf(ctx, "Usage: %s%s", sess.Prefix, mc.usage)
}
