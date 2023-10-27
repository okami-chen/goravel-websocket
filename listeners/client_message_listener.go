package listeners

import (
	"github.com/goravel/framework/contracts/event"
)

type ClientMessageFailListener struct {
}

func (receiver *ClientMessageFailListener) Signature() string {
	return "client_message_fail_listener"
}

func (receiver *ClientMessageFailListener) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     false,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *ClientMessageSuccessListener) Handle(args ...any) error {

	return nil
}

type ClientMessageSuccessListener struct {
}

func (receiver *ClientMessageSuccessListener) Signature() string {
	return "client_message_success_listener"
}

func (receiver *ClientMessageSuccessListener) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *ClientMessageFailListener) Handle(args ...any) error {

	return nil
}
