package bot

import "context"

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
	}

	defaultRegexCommands = []*RegexCommand{
		regTwitter,
	}
}

type Session interface {
	Reply(context.Context, string) error
}

type Command interface {
	Execute(context.Context, Session) error
}
