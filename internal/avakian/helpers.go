package avakian

import (
	"strings"

	"github.com/erei/avakian/internal/bot"
)

func stringInSlice(want string, sl []string) bool {
	for _, s := range sl {
		if strings.EqualFold(want, s) {
			return true
		}
	}

	return false
}

func getBot(s bot.Session) *Bot {
	b, ok := s.Client().(*Bot)
	if !ok {
		panic("non-Bot client returned")
	}

	return b
}
