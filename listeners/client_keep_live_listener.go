package listeners

import (
	"github.com/goravel/framework/contracts/event"
)

type ClientKeepLiveListener struct {
}

func (receiver *ClientKeepLiveListener) Signature() string {
	return "client_keep_live_listener"
}

func (receiver *ClientKeepLiveListener) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *ClientKeepLiveListener) Handle(args ...any) error {

	return nil
}
