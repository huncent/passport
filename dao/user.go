package dao

import (
	"time"

	"github.com/liuhengloveyou/passport/common"
)

type User struct {
	Id         int64      `xorm:"BIGINT(64)"`
	Cellphone  *string    `xorm:"VARCHAR(11)"`
	Email      *string    `xorm:"VARCHAR(45)"`
	Nickname   *string    `xorm:"VARCHAR(45)"`
	Password   *string    `xorm:"not null VARCHAR(45)"`
	AddTime    *time.Time `xorm:"not null TIMESTAMP default 'CURRENT_TIMESTAMP'"`
	UpdateTime *time.Time `xorm:"not null DATETIME updated"`
	Stat       int        `xorm:"not null default 0 INT(11)"`
	Version    int        `xorm:"INT(11) version"`
}

func (p *User) Insert() (e error) {
	_, e = common.Xorms["passport"].InsertOne(p)

	return
}

func (p *User) Update() (e error) {
	_, e = common.Xorms["passport"].Id(p.Id).Update(p)

	return
}

func (p *User) GetOne() (has bool, e error) {
	has, e = common.Xorms["passport"].Get(p)

	return
}
