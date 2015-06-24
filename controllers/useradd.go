package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/models"

	log "github.com/golang/glog"
	"github.com/liuhengloveyou/validator"
)

type UserAdd struct {
	Email string `validate:"noneor,email"`
	Phone string `validate:"noneor,cellphone"`
	Pwd   string
}

func (p *UserAdd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		p.doPost(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (p *UserAdd) doPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		return
	}

	user := &UserAdd{}
	err = json.Unmarshal(body, user)
	if err != nil {
		common.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		return
	}

	if err = validator.Validate(user); err != nil {
		common.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		log.Errorln(*user, err)
		return
	}

	(&models.User{Email: user.Email, Phone: user.Phone, Password: user.Pwd}).Add()

	return
}
