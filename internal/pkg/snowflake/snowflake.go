package snowflake

import (
	"regexp"
	"strconv"
)

var sfRe = regexp.MustCompile(`\d{0,20}`)

func Valid(sf string) bool {
	return sfRe.MatchString(sf)
}

func AsInt64(sf string) (int64, error) {
	i, err := strconv.ParseInt(sf, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}
