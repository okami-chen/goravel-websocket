package events

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

type ClientDisConnectEvent struct {
	UserId   string
	ClientId string
	Time     uint64
}

func NewClientDisConnectEvent(UserId string, ClientId string, Time uint64) error {
	err := facades.Event().Job(&ClientDisConnectEvent{}, []event.Arg{
		{Type: "string", Value: UserId},
		{Type: "string", Value: ClientId},
		{Type: "string", Value: Time},
	}).Dispatch()
	return err
}

func (receiver *ClientDisConnectEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	receiver.UserId = args[0].Value.(string)
	receiver.ClientId = args[1].Value.(string)
	receiver.Time = args[2].Value.(uint64)
	return args, nil
}
