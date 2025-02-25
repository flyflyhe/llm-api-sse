package main

import (
	"bm/internal/config"
	"bm/internal/middleware"
	"bm/internal/tool"
	"bm/pkg/logging"
	"bm/route"
	"flag"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/jinzhu/now"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"runtime/debug"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.Parse()

	now.TimeFormats = append(now.TimeFormats, "200601")

	//初始化
	{
		//日志
		config.InitConfig(configPath)
		logging.InitLogger(logging.Config{
			Debug:     false,
			InfoFile:  "",
			ErrorFile: "",
			CronFile:  "",
		})
		//db.InitDb(config.GetApp().Mysql)
		//db.InitRedis(config.GetApp().Redis)
	}

	engine := server.Default(server.WithHostPorts(fmt.Sprintf("%s:%d", config.GetApp().Web.Ip, config.GetApp().Web.Port)))

	engine.Use(middleware.Print())
	//engine.Use(accesslog.New(accesslog.WithFormat("[${time}] ${status} - ${latency} ${method} ${path} ${queryParams}")))
	//路由
	{
		route.InitPublic(engine.Group("/public"))
	}

	//静态页面
	{
		staticDir := "./static" // 静态文件所在的目录

		// 创建文件服务器
		fileServer := http.FileServer(http.Dir(staticDir))

		// 注册路由，将 / 路径映射到静态文件目录
		http.Handle("/", http.StripPrefix("/", fileServer))

		// 启动 HTTP 服务

		tool.AsyncTask(func() error {
			log.Printf("启动静态文件服务，监听端口 %s，静态文件目录：%s\n", config.GetApp().Web.PortAdmin, staticDir)
			log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.GetApp().Web.IpAdmin, config.GetApp().Web.PortAdmin), nil))
			return nil
		})
	}

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

	engine.Spin()
}

func cronLoop() {
	c := cron.New(cron.WithSeconds())

	logging.Logger.Sugar().Info("cron执行")

	c.Start()
}
