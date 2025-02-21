package api

import (
	"bm/pkg/render"
	"github.com/cloudwego/hertz/pkg/app"
	"time"
)

type Base struct {
}

func (base *Base) Success(c *app.RequestContext, data interface{}) {
	base.SuccessMsg(c, "success", data)
}

func (base *Base) SuccessMsg(c *app.RequestContext, msg string, data interface{}) {
	c.Render(200, render.JSONRender{Data: map[string]interface{}{
		"code":      0,
		"msg":       msg,
		"timestamp": time.Now().Unix(),
		"data":      data,
		"uuid":      c.GetString("uuid"),
	}})
}

func (base *Base) SuccessDefault(c *app.RequestContext) {
	base.SuccessMsg(c, "success", struct{}{})
}

func (base *Base) Fail(c *app.RequestContext, msg string) {
	c.JSON(200, map[string]interface{}{
		"code":      1,
		"msg":       msg,
		"timestamp": time.Now().Unix(),
		"data":      struct{}{},
		"uuid":      c.GetString("uuid"),
	})
}

func Fail(c *app.RequestContext, msg string) {
	c.AbortWithStatusJSON(200, map[string]interface{}{
		"code":      1,
		"msg":       msg,
		"timestamp": time.Now().Unix(),
		"data":      struct{}{},
		"uuid":      c.GetString("uuid"),
	})
}
