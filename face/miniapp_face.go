package face

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/passport/service"
	"github.com/liuhengloveyou/passport/session"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

const LOCATION_WX = "/wx/"

type WxFace struct {
	UserKey    string
	Appid      string
	AppSecrect string
}

func (p *WxFace) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.RequestURI[len(LOCATION_WX):] {
	case "miniapp/login":
		p.miniappLogin(w, r)
		return
	default:
		log.Warningln("404: ", r.RequestURI)
	}

	gocommon.HttpErr(w, http.StatusNotFound, "")

	return
}

func (p *WxFace) miniappLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, "only post")
		return
	}

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Errorln("onlogin read body ERR: ", e.Error())
		gocommon.HttpErr(w, http.StatusBadRequest, "body err.")
		return
	}

	log.Infoln("miniapp.onlogin body: ", string(body))

	app := &service.MiniApp{
		UserKey:    p.UserKey,
		Appid:      p.Appid,
		AppSecrect: p.AppSecrect,
	}
	if e := json.Unmarshal(body, app); e != nil {
		log.Errorln("onlogin Unmarshal body ERR: ", e.Error())
		gocommon.HttpErr(w, http.StatusBadRequest, "body err.")
		return
	}

	info, e := app.Login()
	if e != nil {
		log.Errorln("miniapp user onlogin ERR: ", e.Error())
		gocommon.HttpErr(w, http.StatusInternalServerError, e.Error())
		return
	}

	// wx登录成功
	sess, err := session.GetSession(w, r, "")
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, "会话错误.")
		log.Errorln("session.GetSession ERR:", err.Error())
		return
	}

	sess.Set("user", info)

	w.Write([]byte("{\"sessionid\":\"" + sess.Id("") + "\"}"))

	log.Warningf("login ok: %s %#v", sess.Id(""), info)

	return
}
