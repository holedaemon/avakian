package avakian

import (
	"context"
	"strings"

	"github.com/holedaemon/avakian/internal/bot/message"
)

const (
	discordNameFloor   = 2
	discordNameCeiling = 32
)

var (
	cmdAdminRename = message.NewCommand(
		message.WithCommandArgs(1),
		message.WithCommandFn(cmdAdminRenameFn),
		message.WithCommandPermissions(-1),
	)

	cmdAdmin = message.NewCommand(
		message.WithCommandFn(cmdAdminFn),
		message.WithCommandPermissions(-1),
	)

	adminCommands = message.NewCommandMap(
		message.WithMapScope(true),
		message.WithMapCommand("rename", cmdAdminRename),
		message.WithMapCommand("identitycrisis", cmdAdminRename),
	)
)

func cmdAdminFn(ctx context.Context, s *message.Session) error {
	return adminCommands.ExecuteCommand(ctx, s)
}

func cmdAdminRenameFn(ctx context.Context, s *message.Session) error {
	newName := strings.Join(s.Args, " ")
	if newName == "" {
		return s.Reply(ctx, "Name's empty, dumbass")
	}

	b := getBot(s)
	me := b.Client.State.Me()
	if strings.EqualFold(me.Username, newName) {
		return s.Reply(ctx, "New name is the same as my current")
	}

	if len(newName) < discordNameFloor || len(newName) > discordNameCeiling {
		return s.Replyf(ctx, "Username must be between %d and %d characaters in length", discordNameFloor, discordNameCeiling)
	}

	user := b.Client.User("@me")
	_, err := user.Modify(ctx, newName, "")
	if err != nil {
		return err
	}

	return s.Reply(ctx, "Username has been updated")
}
