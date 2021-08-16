package zapx

import "go.uber.org/zap"

func Command(name string) zap.Field {
	return zap.String("command", name)
}
