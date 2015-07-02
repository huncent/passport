package common

import (
	"encoding/json"
	"os"
)

var (
	ServConfig ConfigServ
)

type ConfigServ struct {
	Listen string `json:"listen"`
	ServID string `json:"serv_id"`
}

func InitServ() {
	initConfig("./app.conf", &ServConfig)
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
