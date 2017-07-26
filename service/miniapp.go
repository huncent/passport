package service

import (
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
	Code       string `json:"code"`
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`

	MiniAppErr
}

func (p *MiniAppUserInfo) Login() error {
	jscode2session := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	_, wxbody, e := gocommon.GetRequest(fmt.Sprintf(jscode2session, common.ServConfig.MiniAppid, common.ServConfig.MiniAppSecrect, p.Code), nil)
	if e != nil {
		log.Errorln("get wx ERR: ", e.Error())
		return e
	}

	log.Infoln("jscode2session res: ", wxbody)
	return nil
}