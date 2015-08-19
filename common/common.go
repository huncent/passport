package common

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/liuhengloveyou/passport/session"
)

var (
	ServConfig ConfigServ
)

type ConfigServ struct {
	Listen  string      `json:"listen"`
	ServID  string      `json:"serv_id"`
	DBs     interface{} `json:"dbs"`
	Session interface{} `json:"session"`
}

func InitPassportServ() error {
	if e := initConfig("./app.conf", &ServConfig); e != nil {
		return e
	}

	if e := InitDbPool(ServConfig.DBs); e != nil {
		return e
	}

	if nil == session.InitDefaultSessionManager(ServConfig.Session) {
		return fmt.Errorf("InitDefaultSessionManager err.")
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
