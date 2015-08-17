package face

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/service"
	"github.com/liuhengloveyou/passport/session"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

func HttpService() {
	http.HandleFunc("/user/add", UserAdd)
	http.HandleFunc("/user/mod", UserModify)
	http.HandleFunc("/user/login", UserLogin)
	http.HandleFunc("/user/auth", UserAuth)
	http.HandleFunc("/user/logout", UserLogout)

	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	s := &http.Server{
		Addr:           common.ServConfig.Listen,
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("GO..." + common.ServConfig.Listen)
	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func UserAdd(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		return

	}
	log.Infoln(string(body))

	user := &service.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	err = user.AddUser()
	if err != nil {
		log.Errorln("user.AddUser() ERR: ", err)
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		return

	}

	gocommon.HttpErr(w, http.StatusOK, nil)
	return
}

func UserModify(w http.ResponseWriter, r *http.Request) {
	//只有登录用户有权修改信息
	sess, err := session.GetSession(w, r)
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		log.Errorln(err.Error())
		return
	}
	if sess.Get("id") == nil {
		gocommon.HttpErr(w, http.StatusForbidden, []byte("用户末登录."))
		log.Warning("update no login:", sess)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	user := &service.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	err = user.UpdateUser()
	if err != nil {
		log.Errorln(*user, err)
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	gocommon.HttpErr(w, http.StatusOK, nil)
	return

}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(w, r)
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		log.Errorln(err.Error())
		return
	}

	user := &service.User{}
	if sess.Get("id") != nil {
		user.Id = sess.Get("id").(int64)
		resp, _ := json.Marshal(user)
		gocommon.HttpErr(w, http.StatusOK, resp)
		log.Warning("login again:", user)
		return // 已经登录
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

	mUser := &service.User{}
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

	has, err := mUser.Get()
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		log.Errorln(*user, err)
		return
	}
	if false == has {
		gocommon.HttpErr(w, http.StatusForbidden, []byte("用户不存在."))
		log.Warningln(user, "用户不存在.")
		return
	}

	loginPWD := service.EncryPWD(mUser.Id, user.Password)
	if loginPWD != mUser.Password {
		gocommon.HttpErr(w, http.StatusForbidden, []byte("用户密码不正确."))
		log.Warningln(*user, *mUser, "用户密码不正确.")
		return

	}

	sess.Set("id", mUser.Id)
	sess.Set("password", mUser.Password)
	log.Infoln(sess)

	gocommon.HttpErr(w, http.StatusOK, []byte("OK"))

	return
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	uid, sid := session.SessionDestroy(w, r)
	log.Infoln(uid, sid)

	if uid != 0 && sid != "" {
		gocommon.HttpErr(w, http.StatusOK, []byte("OK"))
	} else {
		gocommon.HttpErr(w, http.StatusOK, []byte("ERR"))
	}

	return
}

func UserAuth(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(w, r)
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, []byte(err.Error()))
		log.Warningln(err.Error())
		return
	}
	log.Info(sess)

	if sess.Get("id") == nil {
		gocommon.HttpErr(w, http.StatusForbidden, []byte("{}"))
		return
	}

	mUser := &service.User{}
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
