// Package zapx provides useful extentions to Uber's Zap.
package zapx

import (
	"go.uber.org/zap"
)

func New(debug bool) (*zap.Logger, error) {
	if debug {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

		return logger, nil
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func Must(debug bool) *zap.Logger {
	l, err := New(debug)
	if err != nil {
		panic(err)
	}

	return l
}

func Guild(id string) zap.Field {
	return zap.String("guild", id)
}
