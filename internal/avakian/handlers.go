package avakian

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/holedaemon/avakian/internal/bot"
	"github.com/holedaemon/avakian/internal/bot/message"
	"github.com/holedaemon/avakian/internal/bot/reaction"
	"github.com/holedaemon/avakian/internal/bot/regex"
	"github.com/holedaemon/avakian/internal/database/models"
	"github.com/holedaemon/avakian/internal/pkg/modelsx"
	"github.com/holedaemon/avakian/internal/pkg/zapx"
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

	ch, err := b.FetchChannel(ctx, m.ChannelID)
	if err != nil {
		ctxlog.Error(ctx, "error fetching channel", zap.Error(err))
		return
	}

	if ch.Type != discord.ChannelTypeGuildText {
		return
	}

	tx, err := b.DB.BeginTx(ctx, nil)
	if err != nil {
		ctxlog.Error(ctx, "error starting transaction", zap.Error(err))
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

	g, err := modelsx.GetGuildWithPrefixes(ctx, tx, m.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error getting guild", zap.Error(err))
		return
	}

	reSess := regex.NewSession(
		regex.WithSessionClient(b),
		regex.WithSessionGuild(g),
		regex.WithSessionMessage(m),
	)

	if err := regexCommands.ExecuteCommand(ctx, reSess); err != nil {
		switch err {
		case regex.ErrCondition:
			ctxlog.Debug(ctx, "non-fatal condition met in regex command", zap.Error(err))
		case bot.ErrCommandNotExist:
		default:
			ctxlog.Error(ctx, "error running regex command on message", zap.Error(err), zap.String("message", m.Content))
		}
	}

	argv := strings.Split(m.Content, " ")
	prefix := argv[0][:1]

	if prefix != b.DefaultPrefix && !g.HasPrefix(prefix) {
		return
	}

	s := message.NewMessageSession(
		message.WithSessionClient(b),
		message.WithSessionGuild(g),
		message.WithSessionMessage(m),
		message.WithSessionTx(tx),
	)

	ctxlog.Debug(ctx, "attempting to run command")

	err = messageCommands.ExecuteCommand(ctx, s)

	switch err {
	case nil:
		if err := tx.Commit(); err != nil {
			ctxlog.Error(ctx, "error committing transaction", zap.Error(err))
		}
	case bot.ErrCommandNotExist, bot.ErrPermission:
		ctxlog.Debug(ctx, "error returned during command execution", zap.Error(err))
	case bot.ErrUsage:
		if err := messageCommands.SendUsage(ctx, s); err != nil {
			ctxlog.Error(ctx, "error sending command usage", zap.Error(err))
		}
	default:
		ae, ok := err.(*discord.APIError)
		if ok {
			ctxlog.Error(ctx, "http error during command execution", zap.Error(ae), zap.String("message", ae.Message), zap.Int("status", ae.HTTPCode))

			if err := s.Replyf(ctx, "HTTP error encountered during execution: %d %s", ae.HTTPCode, http.StatusText(ae.HTTPCode)); err != nil {
				ctxlog.Error(ctx, "error sending error message", zap.Error(err))
			}
		}

		ctxlog.Error(ctx, "error running command", zap.Error(err))
	}
}

func (b *Bot) handleMessageReactionAdd(r *harmony.MessageReaction) {
	ctx := context.Background()
	ctx = ctxlog.WithLogger(ctx, b.Logger)

	tx, err := b.DB.BeginTx(ctx, nil)
	if err != nil {
		ctxlog.Error(ctx, "error starting transaction", zap.Error(err))
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

	g, err := modelsx.GetGuildWithPrefixes(ctx, tx, r.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error getting guild", zap.Error(err))
		return
	}

	s := reaction.NewSession(
		reaction.WithSessionReaction(r),
		reaction.WithSessionBot(b),
		reaction.WithSessionGuild(g),
		reaction.WithSessionTx(tx),
	)

	err = reactionCommands.ExecuteCommand(ctx, s)
	if err != nil {
		ctxlog.Error(ctx, "error executing command", zap.Error(err))
		return
	}

	if err := tx.Commit(); err != nil {
		ctxlog.Error(ctx, "error committing transaction", zap.Error(err))
	}
}
