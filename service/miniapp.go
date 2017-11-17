package service

import (
	"encoding/json"
	"fmt"

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

	UserId     string `json:"uid"`
	Code       string `json:"code"`
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`

	MiniAppErr
}

type MiniApp struct {
	UserKey    string
	Code       string
	Appid      string
	AppSecrect string
}

func (p *MiniApp) Login() (*MiniAppUserInfo, error) {
	_, wxbody, e := gocommon.GetRequest(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		p.Appid, p.AppSecrect, p.Code), nil)
	if e != nil {
		log.Errorln("jscode2session ERR: ", e.Error())
		return nil, e
	}

	log.Infoln("jscode2session: ", string(wxbody))

	userInfo := &MiniAppUserInfo{}
	if e = json.Unmarshal(wxbody, userInfo); e != nil {
		log.Errorln("jscode2session json ERR: ", e.Error())
		return nil, e
	}

	if userInfo.ErrCode != 0 && userInfo.ErrMsg != "" {
		log.Errorf("jscode2session ERR: %#v", userInfo)
		return nil, fmt.Errorf("jscode2session ERR: %v, %v", userInfo.ErrCode, userInfo.ErrMsg)
	}

	// 生成用户ID
	userInfo.UserId, e = gocommon.AesCBCEncrypt(userInfo.Openid, p.UserKey)
	if e != nil {
		log.Errorln("genuid ERR: ", e.Error())
		return nil, e
	}

	return userInfo, nil
}
