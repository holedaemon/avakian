package bot

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/erei/avakian/internal/database/models"
	"github.com/erei/avakian/internal/pkg/zapx"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/discord"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) handleReady(r *harmony.Ready) {
	b.Logger.Info("connected to Discord")
}

func (b *Bot) handleGuildCreate(g *discord.Guild) {
	b.Logger.Info("received GUILD_CREATE event", zapx.Guild(g.ID))

	ctx := context.Background()
	exists, err := models.Guilds(qm.Where("guild_snowflake = ?", g.ID)).Exists(ctx, b.DB)
	if err != nil {
		b.Logger.Error("error querying db for guild", zapx.Guild(g.ID), zap.Error(err))
		return
	}

	if !exists {
		b.Logger.Info("creating new record for guild", zapx.Guild(g.ID))

		dg := &models.Guild{
			GuildSnowflake: g.ID,
		}

		if err := dg.Insert(ctx, b.DB, boil.Infer()); err != nil {
			b.Logger.Error("error creating db record for guild", zapx.Guild(g.ID), zap.Error(err))
			return
		}

		ctxlog.Info(ctx, "created db record for guild", zapx.Guild(g.ID))
	}
}

func (b *Bot) handleMessage(m *discord.Message) {
	defer func() {
		if r := recover(); r != nil {
			b.Logger.Warn("recovered from panic", zap.Any("recoverer", r))
		}
	}()

	if m.Content == "" || m.Author.Bot {
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	ctx = ctxlog.WithLogger(ctx, b.Logger)

	ctxlog.Debug(ctx, "received message", zap.String("content", m.Content))

	prefs, err := b.GuildPrefixes(ctx, m.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error getting guild prefixes, reverting to default", zap.Error(err))
	}

	argv := strings.Split(m.Content, " ")
	prefix := argv[0][:1]
	cmd := argv[0][1:]

	if prefix != b.DefaultPrefix && !stringInSlice(prefix, prefs) {
		return
	}

	com, ok := defaultMessageCommands[cmd]
	if !ok {
		return
	}

	ch, err := b.FetchChannel(ctx, m.ChannelID)
	if err != nil {
		ctxlog.Error(ctx, "error fetching channel", zap.Error(err))
		return
	}

	g, err := b.FetchGuild(ctx, m.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error fetching guild", zap.Error(err))
		return
	}

	mem, err := b.FetchMember(ctx, m.Author.ID, m.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error fetching member", zap.Error(err))
		return
	}

	p := mem.PermissionsIn(g, ch)
	if !com.HasPermission(p) {
		ctxlog.Debug(ctx, "member does not have required permissions", zap.String("id", m.Author.ID))
		return
	}

	tx, err := b.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		ctxlog.Error(ctx, "error beginning db transaction", zap.Error(err))
		return
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			if errors.Is(err, sql.ErrTxDone) {
				ctxlog.Debug(ctx, "transaction done")
				return
			}

			ctxlog.Error(ctx, "error rolling back transaction", zap.Error(err))
		}
	}()

	sess := b.MessageSession(m)
	sess.Tx = tx

	err = com.Execute(ctx, sess)
	if err != nil {
		ctxlog.Error(ctx, "error during command execution", zap.Error(err))

		switch err := err.(type) {
		case discord.APIError:
			switch err.HTTPCode {
			case http.StatusUnauthorized:
				if err := sess.Reply(ctx, "According to Discord, I'm not authorized to perform whatever it is I'm doing"); err != nil {
					ctxlog.Error(ctx, "error sending message", zap.Error(err))
				}
			case http.StatusForbidden:
				if err := sess.Reply(ctx, "According to Discord, I'm FORBIDDEN from doing whatever I was asked to do"); err != nil {
					ctxlog.Error(ctx, "error sending message", zap.Error(err))
				}
			}
		default:
			if err := sess.Reply(ctx, "An unknown error has occurred, see if you can make sense of it: `"+err.Error()+"`"); err != nil {
				ctxlog.Error(ctx, "error sending message", zap.Error(err))
			}
		}
	} else {
		if err := tx.Commit(); err != nil {
			ctxlog.Error(ctx, "error committing transaction", zap.Error(err))
		}
	}
}
