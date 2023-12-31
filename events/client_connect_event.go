package events

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

type ClientConnectEvent struct {
	UserId   string
	ClientId string
	Time     string
}

func NewClientConnectEvent(UserId string, ClientId string, Time string) error {
	err := facades.Event().Job(&ClientConnectEvent{}, []event.Arg{
		{Type: "string", Value: UserId},
		{Type: "string", Value: ClientId},
		{Type: "string", Value: Time},
	}).Dispatch()
	return err
}

func (receiver *ClientConnectEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	receiver.UserId = args[0].Value.(string)
	receiver.ClientId = args[1].Value.(string)
	receiver.Time = args[2].Value.(string)
	return args, nil
}
