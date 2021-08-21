package snowflake

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

const testStamp = 1590532837033

func TestCanParseTime(t *testing.T) {
	d, err := Time("714971281497522236")
	assert.NilError(t, err, "converting snowflake to time")

	testTime := time.UnixMilli(testStamp)

	assert.Assert(t, d.Equal(testTime), "timestamps are not equal")
}
