package events

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

type ClientKeepLiveEvent struct {
}

func NewClientKeepLiveEvent(UserId string, ClientId string, Time string) error {
	err := facades.Event().Job(&ClientKeepLiveEvent{}, []event.Arg{
		{Type: "string", Value: UserId},
		{Type: "string", Value: ClientId},
		{Type: "string", Value: Time},
	}).Dispatch()
	return err
}

func (receiver *ClientKeepLiveEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
