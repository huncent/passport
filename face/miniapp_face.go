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

const LOCATION_MINIAPP = "/miniapp/"

type miniappFace struct{}

func (p *miniappFace) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.RequestURI[len(LOCATION_MINIAPP):] {
	case "login":
		p.onLogin(w, r)
		return
	default:
		log.Warningln("404: ", r.RequestURI)
	}

	gocommon.HttpErr(w, http.StatusNotFound, "")

	return
}

func (p *miniappFace) onLogin(w http.ResponseWriter, r *http.Request) {
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

	userInfo := &service.MiniAppUserInfo{}
	if e := json.Unmarshal(body, userInfo); e != nil {
		log.Errorln("onlogin Unmarshal body ERR: ", e.Error())
		gocommon.HttpErr(w, http.StatusBadRequest, "body err.")
		return
	}

	log.Infoln("miniapp.onlogin body: ", string(body))

	if e = userInfo.Login(); e != nil {
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

	sess.Set("user", userInfo)

	w.Write([]byte("{\"sessionid\":\"" + sess.Id("") + "\"}"))

	log.Warningf("login ok: %#v", userInfo)

	return
}
