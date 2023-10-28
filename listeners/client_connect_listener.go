package listeners

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

type ClientConnectListener struct {
}

func (receiver *ClientConnectListener) Signature() string {
	return "client_connect_listener"
}

func (receiver *ClientConnectListener) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *ClientConnectListener) Handle(args ...any) error {
	facades.Log().Debug("client_connect_listener")
	return nil
}
