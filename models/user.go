package models

import (
	"fmt"
	"time"

	"github.com/liuhengloveyou/passport/common"
)

type User struct {
	Id         int64     `xorm:"BIGINT(64)"`
	Phone      string    `xorm:"VARCHAR(11)"`
	Email      string    `xorm:"VARCHAR(45)"`
	Password   string    `xorm:"not null VARCHAR(45)"`
	AddTime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	UpdateTime time.Time `xorm:"not null DATETIME updated"`
	Stat       int       `xorm:"not null default 0 INT(11)"`
	Version    int       `xorm:"INT(11) version"`
}

func (p *User) Add() error {
	if e := p.check(); e != nil {
		return e
	}

	r, e := common.Xorms["passport"].Insert(p)

	fmt.Println(r, e)

	return nil
}

/*
 加密用户密码
*/
func (p *User) encryPWD() {

}

/*
 验证字段合法性
*/
func (p *User) check() error {

	return nil
}
