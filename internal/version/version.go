package version

var version string

func Version() string {
	if version == "" {
		return "???"
	}

	return version
}
