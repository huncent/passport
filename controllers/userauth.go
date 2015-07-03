package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/liuhengloveyou/passport/models"
	"github.com/liuhengloveyou/passport/session"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

type UserAuth struct{}

func (p *UserAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		p.doPost(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (p *UserAuth) doPost(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(w, r)
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		log.Warningln(err.Error())
		return
	}
	log.Info(sess)

	mUser := models.UserRequest{}

	if sess.Get("id") != nil {
		mUser.Id = sess.Get("id").(int64)
	}
	if sess.Get("cellphone") != nil {
		mUser.Cellphone = sess.Get("cellphone").(string)
	}
	if sess.Get("email") != nil {
		mUser.Email = sess.Get("email").(string)
	}
	if sess.Get("nickname") != nil {
		mUser.Nickname = sess.Get("nickname").(string)
	}
	log.Infoln(mUser)

	user, _ := json.Marshal(mUser)

	gocommon.HttpErr(w, http.StatusOK, user)

	return
}
