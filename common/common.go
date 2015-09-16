package common

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/liuhengloveyou/passport/session"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	gocommon "github.com/liuhengloveyou/go-common"
)

var (
	ServConfig ConfigServ
	Xorms      = make(map[string]*xorm.Engine)
)

type ConfigServ struct {
	Listen  string      `json:"listen"`
	ServID  string      `json:"serv_id"`
	DBs     interface{} `json:"dbs"`
	Session interface{} `json:"session"`
}

func InitPassportServ(confile string) error {
	if e := gocommon.LoadJsonConfig(confile, &ServConfig); e != nil {
		return e
	}

	if e := InitDbPool(ServConfig.DBs, Xorms); e != nil {
		return e
	}

	if nil == session.InitDefaultSessionManager(ServConfig.Session) {
		return fmt.Errorf("InitDefaultSessionManager err.")
	}

	return nil
}

func InitDbPool(conf interface{}, pool map[string]*xorm.Engine) (err error) {
	var dbs []struct {
		Name   string `json:"name"`
		NGType string `json:"ng_type"`
		DBType string `json:"db_type"`
		DSN    string `json:"dsn"`
	}

	var byteConf []byte
	if byteConf, err = json.Marshal(conf); err != nil {
		return
	}

	if err = json.Unmarshal(byteConf, &dbs); err != nil {
		return
	}

	for _, db := range dbs {
		if db.NGType == "xorm" {
			pool[db.Name], err = xorm.NewEngine(db.DBType, db.DSN)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func InitSystem(conf interface{}) (err error) {
	// 建库,建表
	sqlStr := `
	CREATE DATABASE passport IF NOTEXISTS;

	CREATE TABLE user (
	  id bigint(64) NOT NULL,
	  cellphone varchar(11) COLLATE utf8_bin DEFAULT NULL,
	  email varchar(45) COLLATE utf8_bin DEFAULT NULL,
	  nickname varchar(45) CHARACTER SET utf8 DEFAULT NULL,
	  password varchar(45) COLLATE utf8_bin NOT NULL,
	  add_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  update_time datetime NOT NULL,
	  stat int(11) NOT NULL DEFAULT '0',
	  version int(11) DEFAULT NULL,
	  PRIMARY KEY (id),
	  UNIQUE KEY phone_UNIQUE (cellphone),
	  UNIQUE KEY email_UNIQUE (email),
	  UNIQUE KEY nickname_UNIQUE (nickname)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;`

	var dbs []struct {
		Name   string `json:"name"`
		NGType string `json:"ng_type"`
		DBType string `json:"db_type"`
		DSN    string `json:"dsn"`
	}

	var byteConf []byte
	if byteConf, err = json.Marshal(conf); err != nil {
		return err
	}

	if err := json.Unmarshal(byteConf, &dbs); err != nil {
		return err
	}

	for _, db := range dbs {
		dbc, err := sql.Open(db.DBType, db.DSN)
		if err != nil {
			return err
		}

		if _, err = dbc.Exec(sqlStr); err != nil {
			return err
		}
	}

	return nil
}
