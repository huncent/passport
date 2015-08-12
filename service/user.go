package service

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	p.Id = genUserID()
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

	p.encryPWD()

	return p.toDao().Insert()
}

func (p *User) UpdateUser() (e error) {
	if p.Id <= 0 {
		return fmt.Errorf("user.Id nil")
	}

	if e = validator.Validate(nur); e != nil {
		return
	}

	p.encryPWD()

	return p.toDao().Update()
}

////////

func (p *User) toDao() (dao *dao.User, e error) {
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
		p.Password = encryPWD(p.Id, p.Password)
	}
}

// 加密用户密码
func encryPWD(userid int64, password string) string {
	const SYS_PWD = "When you forgive, You love. And when you love, God's light shines on you."

	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%v%v%v", SYS_PWD, password, (userid/1986)>>4))))
}

// 生成用户ID
func genUserID() int64 {
	id, _ = strconv.ParseInt(gid.ID(), 10, 64)
	return id
}
