package modelsx

import (
	"database/sql"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func makeDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("AVAKIAN_MODELSX_TEST_DSN")
	assert.Assert(t, dsn != "", "dsn empty")

	db, err := sql.Open("pgx", dsn)
	assert.NilError(t, err, "opening db conn")

	err = db.Ping()
	assert.NilError(t, err, "pinging db")

	return db
}
