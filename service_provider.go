package websocket

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/okami-chen/goravel-websocket/events"
	"github.com/okami-chen/goravel-websocket/http/controllers"
	"github.com/okami-chen/goravel-websocket/listeners"
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

	app.MakeEvent().Register(receiver.listen())
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	receiver.Router()
}

func (receiver *ServiceProvider) Router() {
	r := App.MakeRoute()
	r.Prefix("websocket").Get("ws", controllers.NewWebsocketController().Server)
	r.Prefix("api/websocket").Get("test", controllers.NewWebsocketController().Test)
	r.Prefix("api/websocket").Get("register", controllers.NewWebsocketController().Register)
	r.Prefix("api/websocket").Get("send_to_client", controllers.NewWebsocketController().SendToClient)
	go servers.Manager.Start()
	go servers.WriteMessage()
	servers.PingTimer()
}

func (receiver *ServiceProvider) listen() map[event.Event][]event.Listener {
	return map[event.Event][]event.Listener{
		&events.ClientConnectEvent{}: {
			&listeners.ClientConnectListener{},
		},
		&events.ClientDisConnectEvent{}: {
			&listeners.ClientDisConnectListener{},
		},
		&events.ClientKeepLiveEvent{}: {
			&listeners.ClientKeepLiveListener{},
		},
		&events.ClientKillEvent{}: {
			&listeners.ClientKillListener{},
		},
		&events.ClientMessageFailEvent{}: {
			&listeners.ClientMessageFailListener{},
		},
		&events.ClientMessageSuccessEvent{}: {
			&listeners.ClientMessageSuccessListener{},
		},
	}
}
