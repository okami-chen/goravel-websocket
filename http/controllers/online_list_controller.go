package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/okami-chen/goravel-websocket/servers"
)

func (r *WebsocketController) OnelineList(ctx http.Context) http.Response {

	validator, err := ctx.Request().Validate(map[string]string{
		"systemId":  "required",
		"groupName": "required",
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
	groupName := ctx.Request().Input("groupName")

	return ctx.Response().Success().Json(http.Json{
		"code": 0,
		"msg":  "success",
		"data": servers.GetOnlineList(&systemId, &groupName),
	})
}
