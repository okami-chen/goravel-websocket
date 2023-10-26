package facades

import (
	websocket "github.com/okami-chen/goravel-websocket"
	"log"

	"github.com/okami-chen/goravel-websocket/contracts"
)

func Websocket() contracts.Websocket {
	instance, err := websocket.App.Make(websocket.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Websocket)
}
