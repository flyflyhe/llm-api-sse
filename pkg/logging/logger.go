package logging

import (
	"context"
	"go.uber.org/zap"
)

type Config struct {
	Debug     bool
	InfoFile  string
	ErrorFile string
	CronFile  string
}

var Logger *ZLog

func InitLogger(config Config) {
	initZap(config)
	Logger = &ZLog{zapLog}
}

type ZLog struct {
	*zap.Logger
}

func (z *ZLog) WithCtx(ctx context.Context) *zap.SugaredLogger {
	return z.Sugar().With("uuid", ctx.Value("uuid"))
}
