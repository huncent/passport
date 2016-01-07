package service

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/dao"

	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/validator"
)

var gid *gocommon.GlobalID = &gocommon.GlobalID{Expand: common.ServConfig.ServID}

func init() {
	return
}

type User struct {
	Userid    string `validate:"-" json:"id,omitempty"`
	Cellphone string `validate:"noneor,cellphone" json:"cellphone,omitempty"`
	Email     string `validate:"noneor,email" json:"email,omitempty"`
	Nickname  string `validate:"noneor,max=20" json:"nickname,omitempty"`
	Password  string `validate:"nonone,min=6,max=64" json:"password,omitempty"`
}

func (p *User) AddUser() (e error) {
	p.pretreat()

	if p.Userid = genUserID(); p.Userid == "" {
		return fmt.Errorf("用户ID空.")
	}

	p.encryPWD()

	return p.toDao().Insert()
}

func (p *User) UpdateUser() (e error) {
	if p.Userid == "" {
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

	p.Userid = *one.Userid
	if one.Cellphone != nil {
		p.Cellphone = *one.Cellphone
	}
	if one.Email != nil {
		p.Email = *one.Email
	}
	if one.Nickname != nil {
		p.Nickname = *one.Nickname
	}
	if one.Password != nil {
		p.Password = *one.Password
	}

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
}

func (p *User) toDao() *dao.User {
	dao := &dao.User{Userid: &p.Userid}
	if p.Cellphone != "" {
		dao.Cellphone = &p.Cellphone
	}
	if p.Email != "" {
		dao.Email = &p.Email
	}
	if p.Nickname != "" {
		dao.Nickname = &p.Nickname
	}
	if p.Password != "" {
		dao.Password = &p.Password
	}

	return dao
}

func (p *User) encryPWD() {
	if p.Password != "" {
		p.Password = EncryPWD(p.Userid, p.Password)
	}
}

// 加密用户密码
func EncryPWD(userid string, password string) string {
	const SYS_PWD = "When you forgive, You love. And when you love, God's light shines on you."
	iuserid, _ := strconv.ParseInt(userid, 10, 64)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v%v%v", SYS_PWD, password, (iuserid/1986)>>4))))
}

// 生成用户ID
func genUserID() string {
	return gid.ID()
}
