package avakian

import (
	"github.com/holedaemon/avakian/internal/bot/message"
	"github.com/holedaemon/avakian/internal/bot/regex"
)

var (
	messageCommands *message.CommandMap
	regexCommands   *regex.CommandMap
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
		message.WithMapCommand("version", cmdVersion),
		message.WithMapCommand("sf", cmdSnowflake),
		message.WithMapCommand("snowflake", cmdSnowflake),
		message.WithMapCommand("admin", cmdAdmin),
	)

	regexCommands = regex.NewCommandMap(
		regex.WithMapCommand(`https?:\/\/twitter.com\/[a-zA-Z_]{4,15}\/status\/\d{1,20}(?:\?s=\d{0,2})?`, regTwitter),
	)
}
