package regex

import (
	"context"
	"regexp"

	"github.com/holedaemon/avakian/internal/bot"
)

type CommandMap struct {
	comMap     map[string]*Command
	regexCache map[string]*regexp.Regexp
}

type CommandMapOption func(*CommandMap)

func WithMapCommand(pattern string, c *Command) CommandMapOption {
	return func(cm *CommandMap) {
		cm.comMap[pattern] = c
		cm.regexCache[pattern] = regexp.MustCompile(pattern)
	}
}

func NewCommandMap(opts ...CommandMapOption) *CommandMap {
	c := &CommandMap{
		comMap:     make(map[string]*Command),
		regexCache: make(map[string]*regexp.Regexp),
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func (cm *CommandMap) ExecuteCommand(ctx context.Context, sess bot.Session) error {
	s, ok := sess.(*Session)
	if !ok {
		panic("unexpected session type passed to FindAndExecute")
	}

	content := s.Msg.Content

	for _, rgx := range cm.regexCache {
		if rgx.MatchString(content) {
			cmd, ok := cm.comMap[rgx.String()]
			if !ok {
				return bot.ErrCommandNotExist
			}

			s.SetMatch(rgx.FindString(content))

			return cmd.Execute(ctx, s)
		}
	}

	return bot.ErrCommandNotExist
}
