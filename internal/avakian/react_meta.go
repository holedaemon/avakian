package avakian

import (
	"context"

	"github.com/holedaemon/avakian/internal/bot/reaction"
)

var reactMeta = reaction.NewCommand(
	reaction.WithCommandFn(reactMetaFn),
)

func reactMetaFn(ctx context.Context, s *reaction.Session) error {
	return s.Reply(ctx, "Ohhhhhhhhhhhhhhh")
}
