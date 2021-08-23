package snowflake

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	someRandomNumber = 4194304
	DiscordEpoch     = 1420070400000
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

func Time(sf string) (time.Time, error) {
	isf, err := AsInt64(sf)
	if err != nil {
		return time.Time{}, err
	}

	ms := isf/someRandomNumber + DiscordEpoch
	return time.UnixMilli(ms), nil
}

func MarkdownTime(sf string) string {
	t, err := Time(sf)
	if err != nil {
		return "???"
	}

	return fmt.Sprintf("<t:%d>", t.Unix())
}
