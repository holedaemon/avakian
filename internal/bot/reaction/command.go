package reaction

import (
	"context"

	"github.com/holedaemon/avakian/internal/bot"
)

type CommandOption func(*Command)

func WithCommandFn(fn func(context.Context, *Session) error) CommandOption {
	return func(c *Command) {
		c.fn = fn
	}
}

type Command struct {
	fn func(context.Context, *Session) error
}

func NewCommand(opts ...CommandOption) *Command {
	c := new(Command)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Command) Execute(ctx context.Context, sess bot.Session) error {
	s, ok := sess.(*Session)
	if !ok {
		panic("wrong session type passed to command")
	}

	return c.fn(ctx, s)
}

//noop
func (c *Command) Usage(ctx context.Context, sess bot.Session) error {
	return nil
}
