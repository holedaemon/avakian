package bot

import "context"

var defaultMessageCommands map[string]*MessageCommand

func init() {
	defaultMessageCommands = map[string]*MessageCommand{
		"test": cmdTest,
		"flag": cmdFlag,

		"addprefix":    cmdAddPrefix,
		"removeprefix": cmdRemovePrefix,
		"delprefix":    cmdRemovePrefix,
		"listprefix":   cmdListPrefix,
		"listprefixes": cmdListPrefix,
	}
}

type Session interface {
	Reply(context.Context, string) error
}

type Command interface {
	Execute(context.Context, Session) error
}
