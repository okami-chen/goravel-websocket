package websocket

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/okami-chen/goravel-websocket/http/controllers"
	"github.com/okami-chen/goravel-websocket/servers"
)

const Binding = "websocket"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return nil, nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	receiver.Router()
}

func (receiver *ServiceProvider) Router() {
	r := App.MakeRoute()
	r.Prefix("/").Get("ws", controllers.NewWebsocketController().Server)
	r.Prefix("/websocket").Get("test", controllers.NewWebsocketController().Test)
	go servers.Manager.Start()
	go servers.WriteMessage()
	servers.PingTimer()
}
