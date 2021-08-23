package avakian

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/skwair/harmony/discord"
)

type MessageCache struct {
	cache *cache.Cache
}

func NewMessageCache() *MessageCache {
	mc := new(MessageCache)
	mc.cache = cache.New(cache.NoExpiration, time.Hour*1)
	return mc
}

func (mc *MessageCache) Get(id string) (*discord.Message, bool) {
	msg, found := mc.cache.Get(id)
	if found {
		m, ok := msg.(*discord.Message)
		if ok {
			return m, true
		}
	}

	return nil, false
}

func (mc *MessageCache) Set(m *discord.Message) {
	mc.cache.Set(m.ID, m, cache.NoExpiration)
}
