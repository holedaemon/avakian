// Package bot implements Avakian's Discord client.
package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
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

	Twitter *twitter.Client
	DB      *sql.DB
	Logger  *zap.Logger
	Client  *harmony.Client
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

	if b.Twitter == nil {
		fmt.Fprintln(os.Stderr, "[WARN] no twitter client was passed, therefore twitter features will not work")
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

func (b *Bot) QueryGuild(ctx context.Context, sf string) (*models.Guild, error) {
	return models.Guilds(qm.Where("guild_snowflake = ?", sf)).One(ctx, b.DB)
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

func (b *Bot) FetchMemberPermissions(ctx context.Context, gid, cid, mid string) (int, error) {
	gd, err := b.FetchGuild(ctx, gid)
	if err != nil {
		return 0, err
	}

	ch, err := b.FetchChannel(ctx, cid)
	if err != nil {
		return 0, err
	}

	mb, err := b.FetchMember(ctx, mid, gid)
	if err != nil {
		return 0, err
	}

	return mb.PermissionsIn(gd, ch), nil
}

func (b *Bot) MessageSession(msg *discord.Message) *MessageSession {
	argv := strings.Split(msg.Content, " ")
	prefix := argv[0][:1]

	return &MessageSession{
		Msg:    msg,
		Bot:    b,
		Argv:   argv,
		Args:   argv[1:],
		Prefix: prefix,
	}
}

func (b *Bot) RegexSession(msg *discord.Message) *RegexSession {
	return &RegexSession{
		Msg: msg,
		Bot: b,
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

func (b *Bot) AddRole(ctx context.Context, gid, mid, rid, reason string) error {
	g := b.Client.Guild(gid)

	if reason == "" {
		return g.AddMemberRole(ctx, mid, rid)
	}

	return g.AddMemberRoleWithReason(ctx, mid, rid, reason)
}

func (b *Bot) RemoveRole(ctx context.Context, gid, mid, rid, reason string) error {
	g := b.Client.Guild(gid)

	if reason == "" {
		return g.RemoveMemberRole(ctx, mid, rid)
	}

	return g.RemoveMemberRoleWithReason(ctx, mid, rid, reason)
}

func (b *Bot) CreateRole(ctx context.Context, gid, reason string, settings ...discord.RoleSetting) (*discord.Role, error) {
	g := b.Client.Guild(gid)

	roleSettings := discord.NewRoleSettings(settings...)

	if reason == "" {
		return g.NewRole(ctx, roleSettings)
	}

	return g.NewRoleWithReason(ctx, roleSettings, reason)
}

func (b *Bot) DeleteRole(ctx context.Context, gid, rid, reason string) error {
	g := b.Client.Guild(gid)

	if reason == "" {
		return g.DeleteRole(ctx, rid)
	}

	return g.DeleteRoleWithReason(ctx, rid, reason)
}
