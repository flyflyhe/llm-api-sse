package route

import (
	"bm/api"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/hertz-contrib/sse"
	"net/http"
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

	group.POST("/stream", func(ctx context.Context, rc *app.RequestContext) {
		lastEventID := sse.GetLastEventID(rc)
		hlog.CtxInfof(ctx, "last event ID: %s", lastEventID)

		// you must set status code and response headers before first render call
		rc.SetStatusCode(http.StatusOK)
		s := sse.NewStream(rc)
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			event := &sse.Event{
				Event: "timestamp",
				Data:  []byte(fmt.Sprintf("%d", i)),
				ID:    fmt.Sprintf("%d", i),
			}
			err := s.Publish(event)
			if err != nil {
				return
			}
		}

		if err := s.Publish(&sse.Event{Event: "done"}); err != nil {
			return
		}
	})

	group.POST("/ds", ds.ChatCompletion)
	group.GET("/ds-tool", ds.TestFunction)
}
