package bot

import (
	"database/sql"

	"github.com/skwair/harmony"
	"go.uber.org/zap"
)

type Option func(*Bot)

func WithDefaultPrefix(p string) Option {
	return func(b *Bot) {
		b.DefaultPrefix = p
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(b *Bot) {
		b.Logger = l
	}
}

func WithDB(d *sql.DB) Option {
	return func(b *Bot) {
		b.DB = d
	}
}

func WithClient(c *harmony.Client) Option {
	return func(b *Bot) {
		b.Client = c
	}
}
