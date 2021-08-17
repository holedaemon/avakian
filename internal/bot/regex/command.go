package regex

import (
	"context"
	"errors"
)

var ErrCondition = errors.New("regex: a non-fatal condition prevented the command from running")

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

	for _, o := range opts {
		o(c)
	}

	return c
}

func (rg *Command) Execute(ctx context.Context, s *Session) error {
	return rg.fn(ctx, s)
}

// no-op
func (rg *Command) Usage(context.Context, *Session) error {
	return nil
}
