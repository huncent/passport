package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/liuhengloveyou/passport/session"

	redis "github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	gocommon "github.com/liuhengloveyou/go-common"
)

var (
	ServConfig ConfigServ
	DBs        = make(map[string]*gocommon.DBmysql)
	RedisPool  [3]*redis.Pool
)

type ConfigServ struct {
	Listen  string      `json:"listen"`
	ServID  string      `json:"serv_id"`
	Redis   string      `json:"redis"`
	DBs     interface{} `json:"dbs"`
	Session interface{} `json:"session"`

	MiniAppid      string `json:"appid"`
	MiniAppSecrect string `json:"appsecrect"`
	UserKey        string `json:"user_key"`
}

type NilWriter struct{}

func (p *NilWriter) Write(b []byte) (n int, err error) { return 0, nil }

func InitPassportServ(confile string) error {
	if e := gocommon.LoadJsonConfig(confile, &ServConfig); e != nil {
		return e
	}

	if e := gocommon.InitDBPool(ServConfig.DBs, DBs); e != nil {
		return e
	}

	if nil == session.InitDefaultSessionManager(ServConfig.Session) {
		return fmt.Errorf("InitDefaultSessionManager err.")
	}

	// redis 连接池
	for i := 0; i < len(RedisPool); i++ {
		RedisPool[i] = newRedisPool(ServConfig.Redis, i)
	}

	return nil
}

func InitSystem(conf interface{}) (err error) {
	// 建库,建表
	sqlStr := `
	CREATE TABLE user (
	  userid varchar(45) NOT NULL,
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

func newRedisPool(addr string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			conn, e = redis.Dial("tcp", addr, redis.DialDatabase(db))
			if e != nil {
				conn = nil
			}

			return
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) (e error) {
			_, e = c.Do("PING")
			return
		},
	}
}
