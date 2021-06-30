package bot

import (
	"context"
	"net/http"
	"strings"

	"github.com/skwair/harmony"
	"github.com/skwair/harmony/discord"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) handleReady(r *harmony.Ready) {
	b.Logger.Info("connected to Discord")
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

	b.Logger.Debug("received message", zap.String("content", m.Content))

	argv := strings.Split(m.Content, " ")
	prefix := argv[0][:1]
	cmd := argv[0][1:]

	if prefix != b.DefaultPrefix {
		return
	}

	com, ok := defaultMessageCommands[cmd]
	if !ok {
		return
	}

	ctx := context.Background()
	ctx = ctxlog.WithLogger(ctx, b.Logger)

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

	sess := b.MessageSession(m)
	if err := com.Execute(ctx, sess); err != nil {
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

	}
}
