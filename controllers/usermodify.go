package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/passport/models"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/validator"
)

type UserModify struct{}

func (p *UserModify) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		p.doPost(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (p *UserModify) doPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		return
	}

	user := &models.UserRequest{}
	err = json.Unmarshal(body, user)
	if err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		return
	}

	if user.Nickname == "" && user.Password == "" {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte("用户昵称和密码可更新."))
		log.Errorln("usermodify ERR: ", *user)
		return
	}

	if err = validator.Validate(user); err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln(*user, err)
		return
	}

	mUser := &models.User{}
	if user.Nickname != "" {
		mUser.Nickname = user.Nickname
	}
	if user.Password != "" {
		mUser.Password = user.Password
	}

	err = mUser.Update()
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		log.Errorln(*user, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, nil)

}
