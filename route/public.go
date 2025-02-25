package route

import (
	"bm/api"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	"time"
)

type B struct {
	Name string `json:"name"`
}

var (
	ds api.Ds
)

func InitPublic(group *route.RouterGroup) {
	group.GET("/get", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, "111")
		return
	})

	group.POST("/post", func(ctx context.Context, c *app.RequestContext) {
		var reqForm B
		if err := c.Bind(&reqForm); err != nil {
			c.JSON(400, "111")
			return
		}

		c.JSON(200, reqForm)
		return
	})

	group.GET("/stream", func(ctx context.Context, rc *app.RequestContext) {
		rc.Response.Header.Set("Content-Type", "text/event-stream")
		rc.Response.Header.Set("Cache-Control", "no-cache")
		rc.Response.Header.Set("Connection", "keep-alive")

		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			rc.JSON(200, fmt.Sprintf("%d\n", i))
		}
	})

	group.GET("/ds", ds.ChatCompletion)
}
