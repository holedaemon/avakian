package bot

import (
	"context"

	"github.com/erei/avakian/internal/database/models"
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
	QueryGuild(context.Context) (*models.Guild, error)
}

type Command interface {
	Execute(context.Context, Session) error
	Usage(context.Context, Session) error
}
