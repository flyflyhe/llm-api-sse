package logging

import (
	"context"
	"go.uber.org/zap"
)

var Logger *ZLog

func InitLogger(debug bool) {
	initZap(debug)
	Logger = &ZLog{zapLog}
}

type ZLog struct {
	*zap.Logger
}

func (z *ZLog) WithCtx(ctx context.Context) *zap.SugaredLogger {
	return z.Sugar().With("uuid", ctx.Value("uuid"))
}
