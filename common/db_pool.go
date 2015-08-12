package common

import (
	"encoding/json"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	Xorms = make(map[string]*xorm.Engine)
)

func InitDbPool(fn string) error {
	var dbs []struct {
		Name   string `json:"name"`
		NGType string `json:"ng_type"`
		DBType string `json:"db_type"`
		DSN    string `json:"dsn"`
	}

	r, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer r.Close()

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&dbs); err != nil {
		return err
	}

	for _, db := range dbs {
		if db.NGType == "xorm" {
			Xorms[db.Name], err = xorm.NewEngine(db.DBType, db.DSN)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
