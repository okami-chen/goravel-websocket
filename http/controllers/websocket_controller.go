package controllers

import (
	"fmt"
	"github.com/goravel/framework/facades"
	"github.com/okami-chen/goravel-websocket/servers"
	"github.com/okami-chen/goravel-websocket/tools/util"
	nethttp "net/http"

	"github.com/goravel/framework/contracts/http"
	"github.com/gorilla/websocket"
)

type WebsocketController struct {
	// Dependent services
}

func NewWebsocketController() *WebsocketController {
	return &WebsocketController{
		// Inject services
	}
}

func (r *WebsocketController) Server(ctx http.Context) http.Response {
	upGrader := websocket.Upgrader{
		ReadBufferSize:  4096, // Specify the read buffer size
		WriteBufferSize: 4096, // Specify the write buffer size
		// Detect request origin
		CheckOrigin: func(r *nethttp.Request) bool {
			if r.Method != "GET" {
				fmt.Println("method is not GET")
				return false
			}
			if r.URL.Path != "/ws" {
				fmt.Println("path error")
				return false
			}
			return true
		},
	}

	conn, err := upGrader.Upgrade(ctx.Response().Writer(), ctx.Request().Origin(), nil)
	if err != nil {
		return ctx.Response().Success().Json(http.Json{
			"code": 500,
			"msg":  err.Error(),
			"data": nil,
		})
	}
	systemId := ctx.Request().Input("systemId", "default")
	token := ctx.Request().Input("token")

	payload, err := facades.Auth().Parse(ctx, token)
	if err != nil {
		return ctx.Response().Success().Json(http.Json{
			"code": 403,
			"msg":  err.Error(),
			"data": nil,
		})
	}

	clientId := util.GenClientId()
	clientSocket := servers.NewClient(clientId, systemId, conn)
	clientSocket.UserId = payload.Guard

	servers.Manager.AddClient2SystemClient(systemId, clientSocket)

	//读取客户端消息
	clientSocket.Read()

	// 用户连接事件
	servers.Manager.Connect <- clientSocket

	if err = conn.WriteJSON(http.Json{
		"code": 0,
		"msg":  "connect success",
		"data": http.Json{
			"client": clientId,
		},
	}); err != nil {
		conn.Close()
		return nil
	}
	return nil
}

func (r *WebsocketController) Test(ctx http.Context) http.Response {
	msg := "test"
	clients, err := servers.Manager.GetByUserId("80")
	if err != nil {
		return ctx.Response().Success().Json(http.Json{
			"code": 404,
			"msg":  "client not found",
		})
	}
	for _, client := range clients {
		servers.SendMessage2Client(client.ClientId, "1", 0, "test", &msg)
	}

	return ctx.Response().Success().Json(http.Json{
		"code": 0,
		"msg":  "success",
	})
}
