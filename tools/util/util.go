package util

import (
	"fmt"
	"github.com/goravel/framework/facades"
	"github.com/okami-chen/goravel-websocket/tools/crypto"
	"github.com/sony/sonyflake"
	"strconv"
)

func GenUUID() string {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})

	// 生成ID
	id, err := flake.NextID()
	if err != nil {
		fmt.Printf("failed to generate ID: %s\n", err)
		return ""
	}
	return strconv.FormatUint(id, 10)
}

func GenClientId() string {
	raw := []byte(facades.Config().Env("APP_HOST").(string) + ":" + facades.Config().Env("APP_PORT").(string))
	str, err := crypto.Encrypt(raw, []byte("Adba723b7fe06EEH"))
	if err != nil {
		panic(err)
	}

	return str
}

func IsCluster() bool {
	return false
}
