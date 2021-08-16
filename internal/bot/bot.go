package bot

import (
	"context"
	"errors"
)

var (
	ErrCommandNotExist = errors.New("bot: command does not exist")
	ErrUsage           = errors.New("bot: client send usage")
	ErrPermission      = errors.New("bot: insufficient permission")
)

type Client interface {
	SendMessage(context.Context, string, string) error
	CheckPermission(context.Context, Session) (int, error)
}

type Command interface {
	Execute(context.Context, Session) error
	Usage(context.Context, Session) error
}

type Session interface {
	Client() Client
	Reply(context.Context, string) error
	Replyf(context.Context, string, ...interface{}) error
}

type CommandMap interface {
	ExecuteCommand(context.Context, Session) error
}
