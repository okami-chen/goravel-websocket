package listeners

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

type ClientDisConnectListener struct {
}

func (receiver *ClientDisConnectListener) Signature() string {
	return "client_disconnect_listener"
}

func (receiver *ClientDisConnectListener) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *ClientDisConnectListener) Handle(args ...any) error {
	facades.Log().Info("client_connect_listener")
	return nil
}
