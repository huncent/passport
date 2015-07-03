package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/models"
)

type Passport struct {
	ServDomain string
}

func (p *Passport) UserAdd(cellphone, email, nickname, password string) (status int, response []byte, err error) {
	scellphone, semail, snickname := strings.TrimSpace(cellphone), strings.TrimSpace(email), strings.TrimSpace(nickname)
	if scellphone == "" && semail == "" {
		return 0, nil, fmt.Errorf("用户手机号码,邮箱地址不可同时为空.")
	}

	body, _ := json.Marshal(models.UserRequest{Cellphone: scellphone, Email: semail, Nickname: snickname, Password: password})

	status, _, response, err = common.PostRequest(p.ServDomain+"/user/add", body, nil)

	return
}

func (p *Passport) UserLogin(nickname, cellphone, email, password string) (status int, cookies []*http.Cookie, err error) {
	snickname, scellphone, semail := strings.TrimSpace(nickname), strings.TrimSpace(cellphone), strings.TrimSpace(email)
	if snickname == "" && scellphone == "" && semail == "" {
		return 0, nil, fmt.Errorf("用户昵称,用户手机号码,邮箱地址不可同时为空.")
	}

	body, _ := json.Marshal(models.UserRequest{snickname, scellphone, semail, password})

	status, cookies, _, err = common.PostRequest(p.ServDomain+"/user/login", body, nil)

	return
}
