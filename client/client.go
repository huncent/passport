package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/liuhengloveyou/passport/service"

	gocommon "github.com/liuhengloveyou/go-common"
)

type Passport struct {
	ServDomain string
}

func (p *Passport) UserAdd(cellphone, email, nickname, password string) (status int, response []byte, err error) {
	scellphone, semail, snickname := strings.TrimSpace(cellphone), strings.TrimSpace(email), strings.TrimSpace(nickname)
	if scellphone == "" && semail == "" {
		return 0, nil, fmt.Errorf("用户手机号码,邮箱地址不可同时为空.")
	}

	body, _ := json.Marshal(service.User{Cellphone: scellphone, Email: semail, Nickname: snickname, Password: password})
	status, _, response, err = gocommon.PostRequest(p.ServDomain+"/user/add", body, nil)

	return
}

func (p *Passport) UserLogin(nickname, cellphone, email, password string) (status int, cookies []*http.Cookie, err error) {
	snickname, scellphone, semail := strings.TrimSpace(nickname), strings.TrimSpace(cellphone), strings.TrimSpace(email)
	if snickname == "" && scellphone == "" && semail == "" {
		return 0, nil, fmt.Errorf("用户昵称,用户手机号码,邮箱地址不可同时为空.")
	}

	body, _ := json.Marshal(service.User{Cellphone: scellphone, Email: semail, Nickname: snickname, Password: password})
	status, cookies, _, err = gocommon.PostRequest(p.ServDomain+"/user/login", body, nil)

	return
}

func (p *Passport) UserModify(nickname, password string) (status int, cookies []*http.Cookie, err error) {
	snickname, spassword := strings.TrimSpace(nickname), strings.TrimSpace(password)
	if snickname == "" && spassword == "" {
		return 0, nil, fmt.Errorf("用户昵称和密码不可同时为空.")
	}

	body, _ := json.Marshal(service.User{Nickname: snickname, Password: password})
	status, cookies, _, err = gocommon.PostRequest(p.ServDomain+"/user/mod", body, nil)

	return
}
