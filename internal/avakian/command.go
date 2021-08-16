package avakian

import (
	"github.com/erei/avakian/internal/bot/message"
)

var (
	messageCommands *message.CommandMap
	// defaultRegexCommands []*RegexCommand
)

func init() {
	messageCommands = message.NewCommandMap(
		message.WithMapScope(false),
		message.WithMapCommand("test", cmdTest),
		message.WithMapCommand("flag", cmdFlag),
		message.WithMapCommand("prefix", cmdPrefix),
		message.WithMapCommand("pronouns", cmdPronouns),
		message.WithMapCommand("guildctl", cmdGuildctl),
	)

	// defaultRegexCommands = []*RegexCommand{
	// 	regTwitter,
	// }
}
