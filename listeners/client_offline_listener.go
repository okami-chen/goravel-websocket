package listeners

import (
	"github.com/goravel/framework/contracts/event"
)

type ClientOfflineListener struct {
}

func (receiver *ClientOfflineListener) Signature() string {
	return "client_offline_listener"
}

func (receiver *ClientOfflineListener) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *ClientOfflineListener) Handle(args ...any) error {

	return nil
}
