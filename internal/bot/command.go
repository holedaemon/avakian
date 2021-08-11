package bot

import (
	"context"
)

var (
	defaultMessageCommands map[string]*MessageCommand
	defaultRegexCommands   []*RegexCommand
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

type commandMap map[string]Command
