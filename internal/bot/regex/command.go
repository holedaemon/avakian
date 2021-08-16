package regex

import "context"

type Command struct {
	fn func(context.Context, *Session) error
}

func (rg *Command) Execute(ctx context.Context, s *Session) error {
	return rg.fn(ctx, s)
}

// no-op
func (rg *Command) Usage(context.Context, *Session) error {
	return nil
}
