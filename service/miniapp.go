package service

import (
	"encoding/json"
	"fmt"

	"github.com/liuhengloveyou/passport/common"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

type MiniAppErr struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type LoginRst struct {
	ErrMsg string `json:"errMsg"`
	Code   string `json:"code"`
}

type MiniAppUserInfo struct {
	sessionid string
	Code      string `json:"code"`

	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`

	AvatarUrl string `json:"avatarUrl"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Language  string `json:"language"`
	NickName  string `json:"nickName"`
	Province  string `json:"province"`
	Gender    int    `json:"gender"`

	MiniAppErr
}

func (p *MiniAppUserInfo) Login() error {
	jscode2session := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	// test
	p.Openid = "testopenid000"
	p.SessionKey = "testsession000"
	return nil

	_, wxbody, e := gocommon.GetRequest(fmt.Sprintf(jscode2session, common.ServConfig.MiniAppid, common.ServConfig.MiniAppSecrect, p.Code), nil)
	if e != nil {
		log.Errorln("get wx ERR: ", e.Error())
		return e
	}

	log.Infoln("jscode2session res: ", string(wxbody))

	if e = json.Unmarshal(wxbody, p); e != nil {
		log.Errorln("jscode2session json ERR: ", e.Error())
		return e
	}

	log.Infof("jscode2session res: %#v", *p)

	if p.ErrCode != 0 && p.ErrMsg != "" {
		log.Errorf("jscode2session ERR: %#v", *p)
		return fmt.Errorf("jscode2session ERR: %v, %v", p.ErrCode, p.ErrMsg)
	}

	return nil
}
