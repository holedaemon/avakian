package bot

import (
	"github.com/skwair/harmony"
	"go.uber.org/zap"
)

func (b *Bot) handleReady(r *harmony.Ready) {
	b.Logger.Info("connected to Discord")
}

func (b *Bot) handleMessage(m *harmony.Message) {
	b.Logger.Debug("received message", zap.String("content", m.Content))
}
