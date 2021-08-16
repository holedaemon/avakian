package avakian

import (
	"database/sql"

	"github.com/dghubble/go-twitter/twitter"
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

func WithDebug(d bool) Option {
	return func(b *Bot) {
		b.Debug = d
	}
}

func WithTwitter(t *twitter.Client) Option {
	return func(b *Bot) {
		b.Twitter = t
	}
}

func WithAdmins(admins []string) Option {
	return func(b *Bot) {
		b.Admins = admins
	}
}
