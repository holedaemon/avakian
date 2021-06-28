// Package bot implements Avakian's Discord client.
package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/erei/avakian/internal/pkg/zapx"
	"github.com/skwair/harmony"
	"go.uber.org/zap"
)

const (
	defaultPrefix = "!"
)

var (
	ErrClientOption = errors.New("bot: missing required option")
)

type Bot struct {
	Debug         bool
	DefaultPrefix string

	DB     *sql.DB
	Logger *zap.Logger
	Client *harmony.Client
}

func NewBot(opts ...Option) (*Bot, error) {
	b := &Bot{}

	for _, o := range opts {
		o(b)
	}

	if err := b.defaults(); err != nil {
		if !b.Debug {
			return nil, err
		}
	}

	b.Client.OnReady(b.handleReady)
	b.Client.OnMessageCreate(b.handleMessage)

	return b, nil
}

func (b *Bot) defaults() error {
	if b.Client == nil {
		return fmt.Errorf("%w: *harmony.Client", ErrClientOption)
	}

	if b.DB == nil {
		return fmt.Errorf("%w: *sql.DB", ErrClientOption)
	}

	if b.Logger == nil {
		b.Logger = zapx.Must(false)
	}

	if b.DefaultPrefix == "" {
		b.DefaultPrefix = defaultPrefix
	}

	return nil
}

func (b *Bot) Connect(ctx context.Context) error {
	return b.Client.Connect(ctx)
}

func (b *Bot) Disconnect() {
	b.Client.Disconnect()
}
