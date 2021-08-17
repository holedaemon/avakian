package modelsx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/holedaemon/avakian/internal/database/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const pronounQuery = "pronoun = LOWER(?) AND guild_snowflake = ?"

func GetPronoun(ctx context.Context, exec boil.ContextExecutor, sf, prn string) (*models.Pronoun, error) {
	dp, err := models.Pronouns(qm.Where(pronounQuery, prn, sf)).One(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return dp, nil
}

func DeletePronoun(ctx context.Context, exec boil.ContextExecutor, sf, prn string) (string, bool, error) {
	dp, err := GetPronoun(ctx, exec, sf, prn)
	if err != nil {
		return "", false, err
	}

	if dp == nil {
		return "", false, nil
	}

	id := dp.RoleSnowflake

	err = dp.Delete(ctx, exec)
	if err != nil {
		return "", false, err
	}

	return id, true, nil
}

func PronounExists(ctx context.Context, exec boil.ContextExecutor, sf, prn string) (bool, error) {
	dp, err := GetPronoun(ctx, exec, sf, prn)
	if err != nil {
		return false, nil
	}

	return dp != nil, nil
}

func Pronouns(ctx context.Context, exec boil.ContextExecutor, sf string) (models.PronounSlice, error) {
	return models.Pronouns(qm.Where("guild_snowflake = ?", sf)).All(ctx, exec)
}
