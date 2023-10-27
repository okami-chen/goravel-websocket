package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/okami-chen/goravel-websocket/servers"
)

type inputData struct {
	ClientId   string `json:"clientId"`
	UserId     string `json:"userId"`
	SendUserId string `json:"sendUserId"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Data       string `json:"data"`
}

func (r *WebsocketController) SendToClient(ctx http.Context) http.Response {
	//
	var data inputData

	validator, err := ctx.Request().Validate(map[string]string{
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

	data.UserId = ctx.Request().Input("userId")
	data.SendUserId = ctx.Request().Input("sendUserId")
	data.ClientId = ctx.Request().Input("clientId")
	data.Code = ctx.Request().InputInt("code")
	data.Msg = ctx.Request().Input("msg")
	data.Data = ctx.Request().Input("data")

	if data.UserId != "" {
		clients, e := servers.Manager.GetByUserId(data.UserId)
		if e != nil {
			return ctx.Response().Success().Json(http.Json{
				"code": 404,
				"msg":  "client not found",
				"data": []string{},
			})
		}

		var messageId []string
		for _, client := range clients {
			id := servers.SendMessage2Client(client.ClientId, data.SendUserId, data.Code, data.Msg, &data.Data)
			messageId = append(messageId, id)
		}
		return ctx.Response().Success().Json(http.Json{
			"code": 0,
			"msg":  "success",
			"data": http.Json{
				"messages": messageId,
			},
		})
	}

	id := servers.SendMessage2Client(data.ClientId, data.SendUserId, data.Code, data.Msg, &data.Data)

	return ctx.Response().Success().Json(http.Json{
		"code": 0,
		"msg":  "success",
		"data": http.Json{
			"messages": []string{id},
		},
	})
}
