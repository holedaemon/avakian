package zapx

import "go.uber.org/zap"

func Guild(id string) zap.Field {
	return zap.String("guild", id)
}

func Member(id string) zap.Field {
	return zap.String("member", id)
}
