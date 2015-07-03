package models

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
	"time"

	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/passport/common"
)

var globalID *gocommon.GlobalID

type UserRequest struct {
	Id        int64  `validate:"-" json:"id,omitempty"`
	Nickname  string `validate:"noneor,max=20" json:"nickname,omitempty"`
	Cellphone string `validate:"noneor,cellphone" json:"cellphone,omitempty"`
	Email     string `validate:"noneor,email" json:"email,omitempty"`
	Password  string `validate:"nonone,min=6,max=24" json:"password,omitempty"`
}

type User struct {
	Id         int64     `xorm:"BIGINT(64)"`
	Cellphone  string    `xorm:"VARCHAR(11)"`
	Email      string    `xorm:"VARCHAR(45)"`
	Nickname   string    `xorm:"VARCHAR(45)"`
	Password   string    `xorm:"not null VARCHAR(45)"`
	AddTime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	UpdateTime time.Time `xorm:"not null DATETIME updated"`
	Stat       int       `xorm:"not null default 0 INT(11)"`
	Version    int       `xorm:"INT(11) version"`
}

func (p *User) Add() (e error) {
	p.genUserID()
	p.encryPWD()

	if e = p.addCheck(); e != nil {
		return
	}

	_, e = common.Xorms["passport"].InsertOne(p)

	return
}

func (p *User) Update() (e error) {
	p.encryPWD()

	_, e = common.Xorms["passport"].Id(p.Id).Update(p)

	return
}

func (p *User) GetOne() (has bool, e error) {
	has, e = common.Xorms["passport"].Get(p)

	fmt.Println(*p)

	return
}

/*
 生成用户ID
*/
func (p *User) genUserID() {
	if globalID == nil {
		globalID = &gocommon.GlobalID{}
		globalID.Init(common.ServConfig.ServID, "")
	}

	p.Id, _ = strconv.ParseInt(<-globalID.Hole, 10, 64)
}

/*
 加密用户密码
*/
func (p *User) encryPWD() {
	if p.Password != "" {
		p.Password = EncryPWD(p.Id, p.Password)
	}
}

/*
 验证字段合法性
*/
func (p *User) addCheck() error {
	if p.Id <= 0 {
		return fmt.Errorf("user.Id nil")
	}
	if p.Cellphone == "" && p.Email == "" {
		return fmt.Errorf("user.Phone and p.Email nil")
	}
	if p.Password == "" {
		return fmt.Errorf("user.Password nil")
	}

	if p.Cellphone != "" {
		p.Cellphone = strings.ToLower(p.Cellphone)
	}
	if p.Email != "" {
		p.Email = strings.ToLower(p.Email)
	}

	return nil
}

func EncryPWD(userid int64, password string) string {
	const SYS_PWD = "When you forgive, You love. And when you love, God's light shines on you."

	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%v%v%v", SYS_PWD, password, (userid/1986)>>4))))
}
