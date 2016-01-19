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
		return "", err
	}

	body, err := json.Marshal(userinfo)
	if err != nil {
		return "", err
	}

	status, _, response, err := gocommon.PostRequest(p.ServAddr+"/user/add", body, nil, nil)
	if err != nil {
		return "", err
	}

	if status != http.StatusOK {
		return "", fmt.Errorf("%s", response)
	}

	rst := make(map[string]string, 0)
	if err = json.Unmarshal(response, &rst); err != nil {
		return "", err
	}

	return rst["userid"], nil
}

func (p *Passport) UserAuth(token string) (sessionInfo []byte, err error) {
	header := &map[string]string{"TOKEN": token}
	status, _, response, err := gocommon.GetRequest(p.ServAddr+"/user/auth", header)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, nil
	}

	return response, nil
}
