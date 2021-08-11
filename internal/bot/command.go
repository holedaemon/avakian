package bot

import (
	"context"
	"errors"
)

var (
	defaultMessageCommands map[string]*MessageCommand
	defaultRegexCommands   []*RegexCommand

	ErrUsage = errors.New("bot: send usage")
)

func init() {
	defaultMessageCommands = map[string]*MessageCommand{
		"test":     cmdTest,
		"flag":     cmdFlag,
		"prefix":   cmdPrefix,
		"pronouns": cmdPronouns,
		"guildctl": cmdGuildctl,
	}

	defaultRegexCommands = []*RegexCommand{
		regTwitter,
	}
}

type Session interface {
	Reply(context.Context, string) error
	Replyf(context.Context, string, ...interface{}) error
}

type Command interface {
	Execute(context.Context, Session) error
	HasPermission(int) bool
	Usage(context.Context, Session) error
}
