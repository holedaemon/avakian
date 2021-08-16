package modelsx

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"gotest.tools/v3/assert"
)

func TestCanGetGuildWithPrefix(t *testing.T) {
	db := makeDB(t)

	boil.DebugMode = true

	ctx := context.Background()

	_, err := GetGuildWithPrefixes(ctx, db, "gfhghdhfdgf")
	assert.NilError(t, err, "getting guild with prefix")
}
