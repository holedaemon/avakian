// Package avakian implements Avakian's Discord client.
package avakian

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/holedaemon/avakian/internal/bot"
	"github.com/holedaemon/avakian/internal/bot/message"
	"github.com/holedaemon/avakian/internal/database/models"
	"github.com/holedaemon/avakian/internal/pkg/zapx"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/discord"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"
)

const (
	defaultPrefix    = "!"
	maxMessageLength = 2000
)

var (
	ErrClientOption = errors.New("avakian: missing required option")
)

type Bot struct {
	Debug         bool
	DefaultPrefix string
	Admins        []string
	Token         string

	MessageCache *MessageCache
	Twitter      *twitter.Client
	DB           *sql.DB
	Logger       *zap.Logger
	Client       *harmony.Client
	HTTP         *http.Client
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
	b.Client.OnMessageReactionAdd(b.handleMessageReactionAdd)

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

	if b.HTTP == nil {
		b.HTTP = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	if b.MessageCache == nil {
		b.MessageCache = NewMessageCache()
	}

	return nil
}

// SendMessage sends a message to a Discord channel, and implements
// the commands.Dispatcher interface.
func (b *Bot) SendMessage(ctx context.Context, target string, message string) error {
	ch := b.Client.Channel(target)
	_, err := ch.SendMessage(ctx, message)
	return err
}

func (b *Bot) CheckPermission(ctx context.Context, sess bot.Session) (int, error) {
	switch s := sess.(type) {
	case *message.Session:
		if b.IsAdmin(s.Msg.Author.ID) {
			return -1, nil
		}

		return b.FetchMemberPermissions(ctx, s.Msg.GuildID, s.Msg.ChannelID, s.Msg.Author.ID)
	default:
		panic("check permission: unknown session type: not implemented?")
	}
}

func (b *Bot) Connect(ctx context.Context) error {
	return b.Client.Connect(ctx)
}

func (b *Bot) Disconnect() {
	b.Client.Disconnect()
}

func (b *Bot) IsAdmin(sf string) bool {
	return stringInSlice(sf, b.Admins)
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

func (b *Bot) FetchMessage(ctx context.Context, cid string, mid string) (*discord.Message, error) {
	msg, found := b.MessageCache.Get(mid)
	if found {
		return msg, nil
	}

	ch := b.Client.Channel(cid)
	m, err := ch.Message(ctx, mid)
	if err != nil {
		return nil, err
	}

	b.MessageCache.Set(m)

	return m, nil
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
