package avakian

import (
	"context"

	"git.sr.ht/~sircmpwn/getopt"
	"github.com/holedaemon/avakian/internal/bot/message"
	"github.com/holedaemon/avakian/internal/version"
	"github.com/skwair/harmony/discord"
)

var (
	cmdTest = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionSendMessages),
		message.WithCommandFn(cmdTestFn),
	)

	cmdFlag = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionSendMessages),
		message.WithCommandFn(cmdFlagFn),
	)

	cmdVersion = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionSendMessages),
		message.WithCommandFn(cmdVersionFn),
	)
)

func cmdTestFn(ctx context.Context, s *message.Session) error {
	return s.Reply(ctx, "Hello")
}

func cmdVersionFn(ctx context.Context, s *message.Session) error {
	return s.Reply(ctx, version.Version())
}

func cmdFlagFn(ctx context.Context, s *message.Session) error {
	opts, _, err := getopt.Getopts(s.Argv, "b")
	if err != nil {
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
