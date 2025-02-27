package middleware

import (
	"bm/internal/tool"
	"bm/pkg/logging"
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type ReqRes struct {
	Header   []string    `json:"header"`
	Host     string      `json:"host"`
	Method   string      `json:"method"`
	Body     interface{} `json:"body"`
	Response interface{} `json:"response"`
	UUid     string      `json:"u_uid"`
}

func Print() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		uuidStr := uuid.New().String()

		ctx = context.WithValue(ctx, "uuid", uuidStr)

		body := c.Request.Body()
		reqRes := ReqRes{
			Header: strings.Split(string(c.Request.Header.Header()), "\r\n"),
			Host:   string(c.Request.Host()),
			Method: string(c.Request.Method()),
			Body:   string(body),
			UUid:   uuidStr,
		}

		if strings.Contains(c.Request.Header.Get("Content-type"), "application/json") {
			_ = json.Unmarshal(body, &reqRes.Body)
		}

		c.Next(ctx)

		tool.AsyncTask(func() error {
			reqRes.Response = string(c.Response.Body())
			if strings.Contains(c.Response.Header.Get("Content-type"), "application/json") {
				_ = json.Unmarshal(c.Response.Body(), &reqRes.Response)
			}

			logging.Logger.Sugar().Info(tool.ToJson(reqRes))

			return nil
		})
	}
}

func CORSMiddleware() app.HandlerFunc {
	return func(ctx context.Context, rc *app.RequestContext) {
		// 设置允许的来源
		origin := rc.Request.Header.Get("Origin")
		if origin != "" {
			rc.Response.Header.Set("Access-Control-Allow-Origin", origin)
		} else {
			rc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		}

		// 设置允许的请求方法
		rc.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		// 设置允许的请求头
		rc.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 设置是否允许携带凭证
		rc.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求（OPTIONS）
		if string(rc.Request.Method()) == http.MethodOptions {
			rc.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 继续处理请求
		rc.Next(ctx)
	}
}
