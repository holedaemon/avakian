package bot

import (
	"strings"
)

func stringInSlice(want string, sl []string) bool {
	for _, s := range sl {
		if strings.EqualFold(want, s) {
			return true
		}
	}

	return false
}

func buildUsage(command string, commands interface{}) string {
	switch commands := commands.(type) {
	case map[string]*MessageCommand:
		var sb strings.Builder

		sb.WriteString(command + " <")

		i := 1
		for k := range commands {
			if i == len(commands) {
				sb.WriteString(k + ">")
				break
			}

			sb.WriteString(k + "|")

			i++
		}

		return sb.String()
	default:
		panic("invalid type passed to buildUsage()")
	}
}
