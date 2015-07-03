package controllers

import (
	"net/http"

	"github.com/liuhengloveyou/passport/session"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

type UserLogout struct{}

func (p *UserLogout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		p.doPost(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (p *UserLogout) doPost(w http.ResponseWriter, r *http.Request) {
	uid, sid := session.SessionDestroy(w, r)
	if uid != "" && sid != "" {
		gocommon.HttpErr(w, http.StatusOK, []byte("true"))
	} else {
		gocommon.HttpErr(w, http.StatusOK, []byte("false"))
	}

	log.Infoln(uid, sid)

	return
}
