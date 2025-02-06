package route

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
)

type B struct {
	Name string `json:"name"`
}

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
}
