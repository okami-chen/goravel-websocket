package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/okami-chen/goravel-websocket/servers"
)

func (r *WebsocketController) CloseClient(ctx http.Context) http.Response {

	validator, err := ctx.Request().Validate(map[string]string{
		"systemId": "required",
		"clientId": "required",
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
	clientId := ctx.Request().Input("clientId")

	servers.CloseClient(clientId, systemId)

	return ctx.Response().Success().Json(http.Json{
		"code": 0,
		"msg":  "success",
		"data": []string{clientId},
	})
}
