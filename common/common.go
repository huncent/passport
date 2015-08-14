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
	DBs    []byte `json:"dbs"`
}

func InitPassportServ() error {
	if e := initConfig("./app.conf", &ServConfig); e != nil {
		return e
	}

	if e := InitDbPool(ServConfig.DBs); e != nil {
		return e
	}

	return nil
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

	return nil
}
