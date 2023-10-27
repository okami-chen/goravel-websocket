package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/okami-chen/goravel-websocket/servers"
)

func (r *WebsocketController) KickUser(ctx http.Context) http.Response {

	validator, err := ctx.Request().Validate(map[string]string{
		"userId":   "required",
		"systemId": "required",
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
	userId := ctx.Request().Input("userId")

	client, err := servers.Manager.GetByUserId(userId)

	if err != nil {
		return ctx.Response().Success().Json(http.Json{
			"code": 500,
			"msg":  err.Error(),
			"data": []string{},
		})
	}

	var ids []string
	for _, c := range client {
		servers.CloseClient(c.ClientId, systemId)
		ids = append(ids, c.ClientId)
	}

	return ctx.Response().Success().Json(http.Json{
		"code": 0,
		"msg":  "success",
		"data": ids,
	})
}
