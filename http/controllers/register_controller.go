package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"time"
)

func (r *WebsocketController) Register(ctx http.Context) http.Response {

	id := ctx.Request().Input("system_id", "default")

	if !facades.Cache().Has("ws:" + id) {
		facades.Cache().Forever("ws:"+id, http.Json{"system_id": id, "time": time.Now().Unix()})
	}

	return ctx.Response().Success().Json(http.Json{
		"code": 0,
		"msg":  "success",
		"data": []string{},
	})
}
