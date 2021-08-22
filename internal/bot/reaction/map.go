package reaction

import (
	"context"
	"fmt"

	"github.com/holedaemon/avakian/internal/bot"
)

type CommandMapOption func(*CommandMap)

func WithMapCommand(name string, c *Command) CommandMapOption {
	return func(cm *CommandMap) {
		cm.intMap[name] = c
	}
}

type CommandMap struct {
	intMap map[string]*Command
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

func (cm *CommandMap) ExecuteCommand(ctx context.Context, sess bot.Session) error {
	s, ok := sess.(*Session)
	if !ok {
		panic("wrong session type passed to command")
	}

	fmt.Println(s.Reaction.Emoji)

	name := s.Reaction.Emoji.ID
	cmd, ok := cm.intMap[name]
	if !ok {
		return bot.ErrCommandNotExist
	}

	return cmd.Execute(ctx, s)
}

func (cm *CommandMap) SendUsage(context.Context, bot.Session) error {
	return nil
}
