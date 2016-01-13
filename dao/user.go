package dao

import (
	"database/sql"
	"time"

	"github.com/liuhengloveyou/passport/common"
)

type User struct {
	Userid    string    `xorm:"VARCHAR(45)"`
	Cellphone *string   `xorm:"VARCHAR(11)"`
	Email     *string   `xorm:"VARCHAR(45)"`
	Nickname  *string   `xorm:"VARCHAR(45)"`
	Password  *string   `xorm:"not null VARCHAR(256)"`
	AddTime   time.Time `xorm:"not null TIMESTAMP default 'CURRENT_TIMESTAMP'"`
	Version   int       `xorm:"INT(11) version"`
}

func (p *User) Insert() (e error) {
	_, e = common.DBs["passport"].Insert("INSERT INTO user values(?,?,?,?,?,?,?);",
		p.Userid, p.Cellphone, p.Email, p.Nickname, p.Password, p.AddTime, 1)

	return
}

func (p *User) UpdateCellphone() (e error) {
	_, e = common.DBs["passport"].Update("UPDATE user SET cellphone=? WHERE userid=?;", p.Cellphone, p.Userid)

	return
}

func (p *User) UpdateEmail() (e error) {
	_, e = common.DBs["passport"].Update("UPDATE user SET email=? WHERE userid=?;", p.Email, p.Userid)

	return
}

func (p *User) UpdateNickname() (e error) {
	_, e = common.DBs["passport"].Update("UPDATE user SET nickname=? WHERE userid=?;", p.Nickname, p.Userid)

	return
}

func (p *User) UpdatePassword() (e error) {
	_, e = common.DBs["passport"].Update("UPDATE user SET password=? WHERE userid=?;", p.Password, p.Userid)

	return
}

func (p *User) Update() (e error) {
	sqlStr := "UPDATE user SET 1=1"
	if p.Cellphone != nil {
		sqlStr += " and cellphone=" + *p.Cellphone
	}
	if p.Email != nil {
		sqlStr += " and email=" + *p.Email
	}
	if p.Nickname != nil {
		sqlStr += " and nickname=" + *p.Nickname
	}
	if p.Password != nil {
		sqlStr += " and password=" + *p.Password
	}

	sqlStr += " WHERE userid=" + p.Userid

	_, e = common.DBs["passport"].Conn.Exec(sqlStr)

	return
}

func (p *User) QueryByCellphone() (e error) {
	var Cellphone, Email, Nickname, Password sql.NullString
	e = common.DBs["passport"].Conn.QueryRow("SELECT userid, cellphone, email, nickname, password FROM user WHERE cellphone=?;", *p.Cellphone).Scan(&p.Userid, &Cellphone, &Email, &Nickname, &Password)
	if Cellphone.Valid {
		p.Cellphone = &Cellphone.String
	}
	if Email.Valid {
		p.Email = &Email.String
	}
	if Nickname.Valid {
		p.Nickname = &Nickname.String
	}
	if Password.Valid {
		p.Password = &Password.String
	}
	return
}
