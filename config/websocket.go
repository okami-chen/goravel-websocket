package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("websocket", map[string]any{
		"interval": config.Env("WEBSOCKET_INTERVAL", 30),
	})
}
