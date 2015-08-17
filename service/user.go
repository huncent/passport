package service

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"

	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/dao"

	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/validator"
)

var gid *gocommon.GlobalID

func init() {
	gid = &gocommon.GlobalID{ServID: common.ServConfig.ServID}
}

type User struct {
	Id        int64  `validate:"-" json:"id,omitempty"`
	Cellphone string `validate:"noneor,cellphone" json:"cellphone,omitempty"`
	Email     string `validate:"noneor,email" json:"email,omitempty"`
	Nickname  string `validate:"noneor,max=20" json:"nickname,omitempty"`
	Password  string `validate:"nonone,min=6,max=24" json:"password,omitempty"`
}

func (p *User) AddUser() (e error) {
	if p.Cellphone == "" && p.Email == "" {
		return fmt.Errorf("用户手机号和邮箱地址同时为空.")
	}
	if p.Password == "" {
		return fmt.Errorf("用户密码为空.")
	}

	if e = validator.Validate(p); e != nil {
		return
	}

	p.pretreat()

	p.Id = genUserID()
	if p.Id <= 0 {
		return fmt.Errorf("用户ID空.")
	}

	p.encryPWD()

	return p.toDao().Insert()
}

func (p *User) UpdateUser() (e error) {
	if p.Id <= 0 {
		return fmt.Errorf("用户ID空.")
	}

	p.Cellphone, p.Email = "", ""

	if p.Nickname == "" && p.Password == "" {
		return fmt.Errorf("只有用户昵称和密码可更新.")
	}

	p.pretreat()

	if e = validator.Validate(p); e != nil {
		return
	}

	if p.Password != "" {
		p.encryPWD()
	}

	return p.toDao().Update()
}

func (p *User) Get() (has bool, e error) {
	p.pretreat()

	one := p.toDao()
	has, e = one.GetOne()
	if e != nil || has == false {
		return
	}

	p.Id = one.Id
	p.Cellphone = one.Cellphone
	p.Email = one.Email
	p.Nickname = one.Nickname
	p.Password = one.Password

	return true, nil
}

////////
func (p *User) pretreat() {
	if p.Cellphone != "" {
		p.Cellphone = strings.ToLower(p.Cellphone)
	}
	if p.Email != "" {
		p.Email = strings.ToLower(p.Email)
	}
	if p.Nickname != "" {
		p.Nickname = strings.ToLower(p.Nickname)
	}
}

func (p *User) toDao() *dao.User {
	return &dao.User{
		Id:        p.Id,
		Cellphone: p.Cellphone,
		Email:     p.Email,
		Nickname:  p.Nickname,
		Password:  p.Password,
	}
}

func (p *User) encryPWD() {
	if p.Password != "" {
		p.Password = EncryPWD(p.Id, p.Password)
	}
}

// 加密用户密码
func EncryPWD(userid int64, password string) string {
	const SYS_PWD = "When you forgive, You love. And when you love, God's light shines on you."

	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%v%v%v", SYS_PWD, password, (userid/1986)>>4))))
}

// 生成用户ID
func genUserID() int64 {
	id, _ := strconv.ParseInt(gid.ID(), 10, 64)
	return id
}
