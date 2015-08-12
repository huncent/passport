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

func InitPassportServ() error {
	return initConfig("./app.conf", &ServConfig)
}

func initConfig(fn string, config interface{}) error {
	r, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer r.Close()

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(config); err != nil {
		return err
	}
}
