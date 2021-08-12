package modelsx

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"gotest.tools/v3/assert"
)

func TestCanGetGuildWithPrefix(t *testing.T) {
	dsn := os.Getenv("MODELSX_TEST_DSN")

	db, err := sql.Open("pgx", dsn)
	assert.NilError(t, err, "opening db")

	defer db.Close()

	guild, err := GetGuildWithPrefixes(context.Background(), db, "gfhghdhfdgf")
	assert.NilError(t, err, "getting guild")

	t.Log(guild)

}
