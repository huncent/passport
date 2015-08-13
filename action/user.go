package action

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/passport/service"

	log "github.com/golang/glog"
)

func AddUserFromHttp(r *http.Request) (int, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		return http.StatusBadRequest, err
	}

	user := &service.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		return http.StatusBadRequest, err
	}

	err = user.AddUser()
	if err != nil {
		log.Errorln("user.AddUser() ERR: ", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func UserModifyFromHttp(r *http.Request) (int, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		return http.StatusBadRequest, err
	}

	user := &service.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		return http.StatusBadRequest, err
	}

	err = user.UpdateUser()
	if err != nil {
		log.Errorln(*user, err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
