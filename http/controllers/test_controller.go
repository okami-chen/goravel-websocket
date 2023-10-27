package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/okami-chen/goravel-websocket/servers"
)

func (r *WebsocketController) Test(ctx http.Context) http.Response {
	msg := "test"
	clients, err := servers.Manager.GetByUserId("80")
	if err != nil {
		return ctx.Response().Success().Json(http.Json{
			"code": 404,
			"msg":  "client not found",
		})
	}
	var messageId []string
	for _, client := range clients {
		id := servers.SendMessage2Client(client.ClientId, "1", 0, "test", &msg)
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
