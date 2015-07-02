package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/passport/models"
	"github.com/liuhengloveyou/passport/session"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/validator"
)

type UserLogin struct {
	Nickname string `validate:"noneor,max=20"`
	Email    string `validate:"noneor,email"`
	Phone    string `validate:"noneor,cellphone"`
	Password string `validate:"nonone,min=6,max=24"`
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		return
	}

	user := &UserLogin{}
	err = json.Unmarshal(body, user)
	if err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		return
	}

	if err = validator.Validate(user); err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln(*user, err)
		return
	}

	mUser := &models.User{Nickname: user.Nickname, Email: user.Email, Phone: user.Phone}
	has, err := mUser.GetOne()
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		log.Errorln(*user, err)
		return
	}

	if false == has {
		gocommon.HttpErr(w, http.StatusForbidden, []byte("用户不存在."))
		log.Warningln(*user, "用户不存在.")
		return
	}

	loginPWD := models.EncryPWD(mUser.Id, user.Password)
	if loginPWD != mUser.Password {
		gocommon.HttpErr(w, http.StatusForbidden, []byte("用户密码不正确."))
		log.Warningln(*user, *mUser, "用户密码不正确.")
		return

	}

	// session
	sess, err := session.GetSession(w, r)
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		log.Warningln(*user, err.Error())
		return
	}

	sess.Set("uid", mUser.Id)
	sess.Set("cellphone", mUser.Phone)
	sess.Set("email", mUser.Email)
	sess.Set("Nickname", mUser.Nickname)

	fmt.Println(sess)

	gocommon.HttpErr(w, http.StatusOK, nil)

	return
}