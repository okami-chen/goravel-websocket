package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/okami-chen/goravel-websocket/servers"
)

func (r *WebsocketController) BindToGroup(ctx http.Context) http.Response {

	validator, err := ctx.Request().Validate(map[string]string{
		"clientId":  "required",
		"userId":    "required",
		"groupName": "required",
		"systemId":  "required",
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
	userId := ctx.Request().Input("userId")
	groupName := ctx.Request().Input("groupName")
	systemId := ctx.Request().Input("systemId")
	if userId != "" {
		clients, e := servers.Manager.GetByUserId(userId)
		if e != nil {
			return ctx.Response().Success().Json(http.Json{
				"code": 404,
				"msg":  "client not found",
				"data": []string{},
			})
		}

		for _, client := range clients {
			servers.AddClient2Group(systemId, groupName, client.ClientId, client.UserId, "")
		}
	} else {
		client, err := servers.Manager.GetByClientId(ctx.Request().Input("clientId"))
		if err != nil {
			return ctx.Response().Success().Json(http.Json{
				"code": 404,
				"msg":  "client not found",
				"data": []string{},
			})
		}
		servers.AddClient2Group(systemId, groupName, client.ClientId, client.UserId, "")
	}

	return ctx.Response().Success().Json(http.Json{
		"code": 0,
		"msg":  "success",
		"data": []string{},
	})
}
