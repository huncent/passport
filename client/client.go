package client

import (
	"net/http"

	gocommon "github.com/liuhengloveyou/go-common"
)

type Passport struct {
	ServAddr string
}

func (p *Passport) UserAdd(data []byte) (status int, response []byte, err error) {
	status, _, response, err = gocommon.PostRequest(p.ServAddr+"/user/add", data, nil)
	return
}

func (p *Passport) UserLogin(data []byte) (status int, cookies []*http.Cookie, err error) {
	status, cookies, _, err = gocommon.PostRequest(p.ServAddr+"/user/login", data, nil)
	return
}

func (p *Passport) UserModify(data []byte) (status int, cookies []*http.Cookie, err error) {
	status, cookies, _, err = gocommon.PostRequest(p.ServAddr+"/user/mod", data, nil)
	return
}
