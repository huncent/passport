package dao

import (
	"time"

	"github.com/liuhengloveyou/passport/common"
)

type User struct {
	Userid     *string    `xorm:"VARCHAR(45)"`
	Cellphone  *string    `xorm:"VARCHAR(11)"`
	Email      *string    `xorm:"VARCHAR(45)"`
	Nickname   *string    `xorm:"VARCHAR(45)"`
	Password   *string    `xorm:"not null VARCHAR(256)"`
	AddTime    *time.Time `xorm:"not null TIMESTAMP default 'CURRENT_TIMESTAMP'"`
	UpdateTime *time.Time `xorm:"not null DATETIME updated"`
	Stat       int        `xorm:"not null default 0 INT(11)"`
	Version    int        `xorm:"INT(11) version"`
}

func (p *User) Insert() (e error) {
	_, e = common.DBs["passport"].Insert("INSERT INTO user values(?,?,?,?,?,?,?,?,?)",
		p.Userid, p.Cellphone, p.Email, p.Nickname, p.Password, p.AddTime, p.UpdateTime, p.Stat, 1)

	return
}

func (p *User) Update() (e error) {
	//	_, e = common.DBs["passport"].Update("UPDATE user"

	return
}

func (p *User) GetOne() (has bool, e error) {
	//	has, e = common.DBs["passport"].Query(sqlStr string, args ...interface{})
	return
}
