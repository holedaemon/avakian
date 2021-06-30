package bot

import (
	"context"

	"git.sr.ht/~sircmpwn/getopt"
	"github.com/skwair/harmony/discord"
)

var cmdTest = &MessageCommand{
	permissions: discord.PermissionSendMessages,
	fn:          cmdTestFn,
}

func cmdTestFn(ctx context.Context, s *MessageSession) error {
	return s.Reply(ctx, "Hello")
}

var cmdFlag = &MessageCommand{
	permissions: discord.PermissionSendMessages,
	fn:          cmdFlagFn,
}

func cmdFlagFn(ctx context.Context, s *MessageSession) error {
	opts, _, err := getopt.Getopts(s.Argv, "b")
	if err != nil {
		s.Reply(ctx, "Unable to parse flags")
		return err
	}

	passedFlag := false
	for _, opt := range opts {
		switch opt.Option {
		case 'b':
			passedFlag = true
		}
	}

	if passedFlag {
		return s.Reply(ctx, "bada bing bada boom ayyyy")
	}

	return s.Reply(ctx, "Flag -b was not passed")
}
