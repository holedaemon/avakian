package message

import (
	"context"
	"strings"

	"github.com/holedaemon/avakian/internal/bot"
)

type CommandMapOption func(*CommandMap)

func WithMapCommand(name string, c *Command) CommandMapOption {
	return func(cm *CommandMap) {
		cm.intMap[name] = c
	}
}

func WithMapScope(scope bool) CommandMapOption {
	return func(cm *CommandMap) {
		cm.scoped = scope
	}
}

func NewCommandMap(opts ...CommandMapOption) *CommandMap {
	cm := &CommandMap{
		intMap: make(map[string]*Command),
	}

	for _, o := range opts {
		o(cm)
	}

	return cm
}

type CommandMap struct {
	intMap map[string]*Command
	scoped bool
}

func (cm *CommandMap) BuildUsage(command string) string {
	var sb strings.Builder

	sb.WriteString(command + " <")

	i := 1
	for k := range cm.intMap {
		if i == len(cm.intMap) {
			sb.WriteString(k + ">")
			break
		}

		sb.WriteString(k + "|")

		i++
	}

	return sb.String()
}

func (cm *CommandMap) ExecuteCommand(ctx context.Context, s bot.Session) error {
	sess, ok := s.(*Session)
	if !ok {
		panic("wrong session type passed to command")
	}

	if len(sess.Argv) == 0 {
		return bot.ErrCommandNotExist
	}

	if cm.scoped {
		name := sess.Args[0]
		cmd, ok := cm.intMap[name]
		if !ok {
			return bot.ErrUsage
		}

		sess.Args = sess.Args[1:]

		if len(sess.Args) < cmd.minArgs {
			return bot.ErrUsage
		}

		return cmd.Execute(ctx, s)
	}

	name := sess.Argv[0][1:]
	cmd, ok := cm.intMap[name]
	if !ok {
		return bot.ErrCommandNotExist
	}

	if len(sess.Args) < cmd.minArgs {
		return bot.ErrUsage
	}

	return cmd.Execute(ctx, sess)
}

func (cm *CommandMap) SendUsage(ctx context.Context, s bot.Session) error {
	sess, ok := s.(*Session)
	if !ok {
		panic("wrong session type passed to command")
	}

	if cm.scoped {
		sess.Args = sess.Args[1:]
	}

	name := sess.Argv[0][1:]
	cmd, ok := cm.intMap[name]
	if !ok {
		return bot.ErrCommandNotExist
	}

	return cmd.Usage(ctx, s)
}
