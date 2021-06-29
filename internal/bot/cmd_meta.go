package bot

import (
	"context"

	"github.com/skwair/harmony/discord"
)

var cmdTest = &MessageCommand{
	permissions: discord.PermissionSendMessages,
	flags:       "",
	fn:          cmdTestFn,
}

func cmdTestFn(ctx context.Context, s *MessageSession) error {
	return s.Reply(ctx, "Hello")
}
