package avakian

import (
	"fmt"
	"strings"

	"github.com/holedaemon/avakian/internal/bot"
	"github.com/holedaemon/avakian/internal/bot/reaction"
	"github.com/skwair/harmony/discord"
)

const jumpLinkURL = "https://discord.com/channels/%s/%s/%s"

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

func fullUsername(u *discord.GuildMember) string {
	if u.Nick != "" {
		return u.Nick
	}

	return u.User.Username
}

func jumpLinkFromSession(s bot.Session) string {
	switch s := s.(type) {
	case *reaction.Session:
		return fmt.Sprintf(jumpLinkURL, s.Reaction.GuildID, s.Reaction.ChannelID, s.Reaction.MessageID)
	}

	return ""
}
