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

	UserId     string `json:"uid"`
	Code       string `json:"code"`
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`

	MiniAppErr
}

func (p *MiniAppUserInfo) Login() error {
	_, wxbody, e := gocommon.GetRequest(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		common.ServConfig.MiniAppid, common.ServConfig.MiniAppSecrect, p.Code), nil)
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

	// 生成用户ID
	p.UserId, e = gocommon.AesCBCEncrypt(p.Openid, common.ServConfig.UserKey)
	if e != nil {
		log.Errorln("genuid ERR: ", e.Error())
		return e
	}

	return nil
}
