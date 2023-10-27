package events

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

type ClientDisConnectEvent struct {
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
	return args, nil
}
