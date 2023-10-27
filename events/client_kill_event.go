package events

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

type ClientKillEvent struct {
}

func NewClientKillEvent(UserId string, ClientId string, Time string) error {
	err := facades.Event().Job(&ClientKillEvent{}, []event.Arg{
		{Type: "string", Value: UserId},
		{Type: "string", Value: ClientId},
		{Type: "string", Value: Time},
	}).Dispatch()
	return err
}

func (receiver *ClientKillEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
