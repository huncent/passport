package common

import (
	"encoding/json"
	"os"
)

var (
	Config ConfigS // 系统配置信息
)

type ConfigS struct {
	Listen string `json:"listen"`
}

func InitServ() {
	initConfig("./app.conf", &Config)
}

func initConfig(fn string, config interface{}) {
	r, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(config); err != nil {
		panic(err)
	}

	r.Close()
}
