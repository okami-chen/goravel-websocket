package events

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

type ClientMessageSuccessEvent struct {
}
type ClientMessageFailEvent struct {
}

func NewClientMessageSuccessEvent(userId, clientId, time, messageId string) error {
	err := facades.Event().Job(&ClientMessageSuccessEvent{}, []event.Arg{
		{Type: "string", Value: userId},
		{Type: "string", Value: clientId},
		{Type: "string", Value: time},
		{Type: "string", Value: messageId},
	}).Dispatch()
	return err
}

func (receiver *ClientMessageSuccessEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}

func NewClientMessageFailEvent(userId, clientId, time, messageId string) error {
	err := facades.Event().Job(&ClientMessageFailEvent{}, []event.Arg{
		{Type: "string", Value: userId},
		{Type: "string", Value: clientId},
		{Type: "string", Value: time},
		{Type: "string", Value: messageId},
	}).Dispatch()
	return err
}

func (receiver *ClientMessageFailEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
