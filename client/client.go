package client

import (
	"net/http"

	gocommon "github.com/liuhengloveyou/go-common"
)

type Passport struct {
	ServAddr string
}

func (p *Passport) Execute(uri string, data []byte, cookies []*http.Cookie) (status int, responseCookies []*http.Cookie, response []byte, err error) {
	status, responseCookies, response, err = gocommon.PostRest(p.ServAddr+uri, data, cookies, nil)
	return
}
