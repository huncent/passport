package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/passport/models"
	"github.com/liuhengloveyou/passport/session"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/validator"
)

type UserLogin struct {
}

func (p *UserLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		p.doPost(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (p *UserLogin) doPost(w http.ResponseWriter, r *http.Request) {
	//
	sess, err := session.GetSession(w, r)
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		log.Errorln(err.Error())
		return
	}

	user := models.UserRequest{}
	if sess.Get("id") != nil {
		user.Id = sess.Get("id").(int64)
		resp, _ := json.Marshal(user)
		gocommon.HttpErr(w, http.StatusOK, resp)
		log.Warning("login again:", user)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		return
	}

	log.Infoln(string(body))

	err = json.Unmarshal(body, &user)
	if err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		return
	}

	log.Infoln(user)

	if err = validator.Validate(user); err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln(user, err)
		return
	}

	//
	mUser := &models.User{}
	if user.Cellphone != "" {
		mUser.Cellphone = user.Cellphone
		sess.Set("cellphone", user.Cellphone)
	} else if user.Email != "" {
		mUser.Email = user.Email
		sess.Set("email", user.Email)
	} else if user.Nickname != "" {
		mUser.Nickname = user.Nickname
		sess.Set("nickname", user.Nickname)
	} else {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte("用户标识为空."))
		log.Errorln("用户标识为空.")
		return
	}

	log.Infoln(mUser)

	has, err := mUser.GetOne()
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		log.Errorln(user, err)
		return
	}

	if false == has {
		gocommon.HttpErr(w, http.StatusForbidden, []byte("用户不存在."))
		log.Warningln(user, "用户不存在.")
		return
	}

	loginPWD := models.EncryPWD(mUser.Id, user.Password)
	if loginPWD != mUser.Password {
		gocommon.HttpErr(w, http.StatusForbidden, []byte("用户密码不正确."))
		log.Warningln(user, *mUser, "用户密码不正确.")
		return

	}

	sess.Set("id", mUser.Id)
	sess.Set("password", mUser.Password)
	sess.Set("version", mUser.Version)
	log.Infoln(sess)

	gocommon.HttpErr(w, http.StatusOK, nil)

	return
}
