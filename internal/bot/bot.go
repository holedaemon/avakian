// Package bot implements Avakian's Discord client.
package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/erei/avakian/internal/database/models"
	"github.com/erei/avakian/internal/pkg/zapx"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/discord"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	b.Client.OnGuildCreate(b.handleGuildCreate)

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

func (b *Bot) FetchGuild(ctx context.Context, id string) (*discord.Guild, error) {
	sg := b.Client.State.Guild(id)
	if sg != nil {
		return sg, nil
	}

	gr := b.Client.Guild(id)
	ag, err := gr.Get(ctx)
	if err != nil {
		return nil, err
	}

	return ag, nil
}

func (b *Bot) FetchChannel(ctx context.Context, id string) (*discord.Channel, error) {
	sc := b.Client.State.Channel(id)
	if sc != nil {
		return sc, nil
	}

	cr := b.Client.Channel(id)
	ac, err := cr.Get(ctx)
	if err != nil {
		return nil, err
	}

	return ac, nil
}

func (b *Bot) FetchMember(ctx context.Context, mid, gid string) (*discord.GuildMember, error) {
	g := b.Client.Guild(gid)

	m, err := g.Member(ctx, mid)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (b *Bot) MessageSession(msg *discord.Message) *MessageSession {
	argv := strings.Split(msg.Content, " ")

	return &MessageSession{
		Msg:  msg,
		Bot:  b,
		Argv: argv,
		Args: argv[1:],
	}
}

func (b *Bot) GuildPrefixes(ctx context.Context, id string) ([]string, error) {
	prefixes, err := models.Prefixes(qm.Where("guild_snowflake = ?", id)).All(ctx, b.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	prefs := make([]string, 0, len(prefixes))

	for _, p := range prefixes {
		prefs = append(prefs, p.Prefix)
	}

	return prefs, nil
}
