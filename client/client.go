package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/liuhengloveyou/passport/service"

	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/validator"
)

type Passport struct {
	ServAddr string
}

func (p *Passport) UserAdd(cellphone, email, nickname, password string) (userid string, err error) {
	userinfo := &service.User{Cellphone: cellphone, Email: email, Nickname: nickname, Password: password}
	if err = validator.Validate(userinfo); err != nil {
		return 0, err
	}

	body, err := json.Marshal(userinfo)
	if err != nil {
		return 0, err
	}

	status, _, response, err := gocommon.PostRequest(p.ServAddr+"/user/add", body, nil, nil)
	if err != nil {
		return 0, err
	}

	if status != http.StatusOK {
		return 0, fmt.Errorf("%s", response)
	}

	rst := make(map[string]int64, 0)
	if err = json.Unmarshal(response, &rst); err != nil {
		return 0, err
	}

	return rst["userid"], nil
}

func (p *Passport) UserAuth(cookies []*http.Cookie) (status int, response []byte, err error) {
	//status, _, response, err = gocommon.PostRequest(p.ServAddr+"/user/auth", make([]byte, 0), cookies, nil)
	return
}

func (p *Passport) Execute(uri string, data []byte, cookies []*http.Cookie) (status int, responseCookies []*http.Cookie, response []byte, err error) {
	//status, responseCookies, response, err = gocommon.PostRequest(p.ServAddr+uri, data, cookies, nil)
	return
}
