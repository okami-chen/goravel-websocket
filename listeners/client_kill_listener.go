package listeners

import (
	"github.com/goravel/framework/contracts/event"
)

type ClientKillListener struct {
}

func (receiver *ClientKillListener) Signature() string {
	return "client_kill_listener"
}

func (receiver *ClientKillListener) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *ClientKillListener) Handle(args ...any) error {

	return nil
}
