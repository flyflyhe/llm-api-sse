package main

import (
	"bm/internal/config"
	"bm/pkg/logging"
	"flag"
	"github.com/jinzhu/now"
	"github.com/robfig/cron/v3"
	"runtime/debug"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.Parse()

	now.TimeFormats = append(now.TimeFormats, "200601")

	//日志
	config.InitConfig(configPath)
	logging.InitLogger(config.GetApp().Web.Debug)
	//db.InitDb(config.GetApp().Mysql)
	//db.InitRedis(config.GetApp().Redis)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logging.Logger.Sugar().Error(string(debug.Stack()))
				logging.Logger.Sugar().Error(err)
			}
		}()

		if config.GetApp().Cron {
			logging.Logger.Sugar().Info("cron开启")
			cronLoop()
		} else {
			logging.Logger.Sugar().Info("cron未开启")
		}
	}()
}

func cronLoop() {
	c := cron.New(cron.WithSeconds())

	logging.Logger.Sugar().Info("cron执行")

	c.Start()
}
