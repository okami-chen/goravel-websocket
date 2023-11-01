package controllers

import (
	"encoding/json"
	"github.com/goravel/framework/contracts/http"
	"github.com/okami-chen/goravel-websocket/servers"
)

func (r *WebsocketController) SendToSystem(ctx http.Context) http.Response {

	validator, err := ctx.Request().Validate(map[string]string{
		"systemId":   "required",
		"sendUserId": "required",
		"code":       "required",
		"msg":        "required",
		"data":       "required",
	})

	if err != nil {
		return ctx.Response().Success().Json(http.Json{
			"code": 500,
			"msg":  err.Error(),
			"data": []string{},
		})
	}

	if validator.Fails() {
		return ctx.Response().Success().Json(http.Json{
			"code": 500,
			"msg":  validator.Errors().One(),
			"data": []string{},
		})
	}

	systemId := ctx.Request().Input("systemId")
	sendUserId := ctx.Request().Input("sendUserId")
	code := ctx.Request().InputInt("code")
	msg := ctx.Request().Input("msg")

	str := ctx.Request().InputMap("data")
	bt, e := json.Marshal(str)
	if e != nil {
		return ctx.Response().Success().Json(http.Json{
			"code": 500,
			"msg":  e.Error(),
			"data": []string{},
		})
	}
	m := string(bt)
	servers.SendMessage2System(systemId, sendUserId, code, msg, m)

	return ctx.Response().Success().Json(http.Json{
		"code": 0,
		"msg":  "success",
		"data": []string{},
	})
}
