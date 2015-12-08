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
	"github.com/liuhengloveyou/validator"
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

func optionsFilter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://web.xim.com:9000")
	w.Header().Add("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "X-API, X-REQUEST-ID, X-API-TRANSACTION, X-API-TRANSACTION-TIMEOUT, X-RANGE, Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)

	return
}

func UserAdd(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		return

	}
	log.Infoln(string(body))

	user := &service.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = validator.Validate(user); err != nil {
		log.Errorln("validator ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		return
	}

	if user.Cellphone == "" && user.Email == "" {
		log.Errorln("ERR: 用户手机号和邮箱地址同时为空.")
		gocommon.HttpErr(w, http.StatusBadRequest, "用户手机号和邮箱地址同时为空.")
		return
	}
	if user.Password == "" {
		log.Errorln("ERR: 用户密码为空.")
		gocommon.HttpErr(w, http.StatusBadRequest, "用户密码为空.")
		return
	}

	err = user.AddUser()
	if err != nil {
		log.Errorln("user.AddUser() ERR: ", err)
		gocommon.HttpErr(w, http.StatusInternalServerError, "服务错误.")
		return

	}

	gocommon.HttpErr(w, http.StatusOK, "OK")
	return
}

func UserModify(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	//只有登录用户有权修改信息
	sess, err := session.GetSession(w, r, nil)
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, err.Error())
		log.Errorln(err.Error())
		return
	}
	if sess.Get("id") == nil {
		gocommon.HttpErr(w, http.StatusForbidden, "用户末登录.")
		log.Warning("update no login:", sess)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		return
	}

	user := &service.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		return
	}

	err = user.UpdateUser()
	if err != nil {
		log.Errorln(*user, err)
		gocommon.HttpErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, "OK")
	return

}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	sess, err := session.GetSession(w, r, nil)
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, err.Error())
		log.Errorln(err.Error())
		return
	}

	user := &service.User{}
	tmp := sess.Get("user")
	if tmp != nil {
		user = tmp.(*service.User)
		resp, _ := json.Marshal(user)
		gocommon.HttpErr(w, http.StatusOK, string(resp))
		log.Warning("login again:", user)
		return // 已经登录
	}
	// 新登录。。。

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		return
	}
	log.Infoln(string(body))

	if err := json.Unmarshal(body, user); err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		return
	}

	if user.Cellphone == "" && user.Email == "" && user.Nickname != "" {
		gocommon.HttpErr(w, http.StatusBadRequest, "用户标识为空.")
		log.Errorln("用户标识为空.")
		return
	}

	if err = validator.Validate(user); err != nil {
		log.Errorln("validator ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		return
	}

	mUser := &service.User{Cellphone: user.Cellphone, Email: user.Email, Nickname: user.Nickname}
	has, err := mUser.Get()
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, err.Error())
		log.Errorln(*user, err)
		return
	}
	if false == has {
		gocommon.HttpErr(w, http.StatusForbidden, "用户不存在.")
		log.Warningln(user, "用户不存在.")
		return
	}

	loginPWD := service.EncryPWD(mUser.Id, user.Password)
	if loginPWD != mUser.Password {
		gocommon.HttpErr(w, http.StatusForbidden, "用户密码不正确.")
		log.Warningln(*user, *mUser, "用户密码不正确.")
		return

	}

	sess.Set("user", user)
	log.Infoln(sess)

	gocommon.HttpErr(w, http.StatusOK, "OK")

	return
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	sid := session.SessionDestroy(w, r)
	log.Infoln(sid)

	if sid != "" {
		gocommon.HttpErr(w, http.StatusOK, "OK")
	} else {
		gocommon.HttpErr(w, http.StatusOK, "ERR")
	}

	return
}

func UserAuth(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	sess, err := session.GetSession(w, r, nil)
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, err.Error())
		log.Warningln(err.Error())
		return
	}
	log.Info(sess)

	if sess.Get("id") == nil {
		gocommon.HttpErr(w, http.StatusForbidden, "{}")
		log.Errorln("session no id:", sess)
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

	gocommon.HttpErr(w, http.StatusOK, string(user))

	return
}
