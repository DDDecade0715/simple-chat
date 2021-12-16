package initialize

import (
	zap2 "gin-derived/pkg/zap"
	"go.uber.org/zap"
)

func InitLogger(s string, p string) *zap.SugaredLogger {
	return zap2.GetZap(s, p)
}
