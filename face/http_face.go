package face

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

	http.Handle(LOCATION_WX, &WxFace{})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("404: ", r.RequestURI)
		w.WriteHeader(http.StatusNotFound)
	})

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

/*
 * 跨域资源共享
 */
func optionsFilter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://web.xim.com:9000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "X-API, X-REQUEST-ID, X-API-TRANSACTION, X-API-TRANSACTION-TIMEOUT, X-RANGE, Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)

	return
}

func AuthFilter(w http.ResponseWriter, r *http.Request) (sess session.SessionStore, auth bool) {
	token := strings.TrimSpace(r.Header.Get("sessionid"))
	if token == "" {
		if cookie, e := r.Cookie("sessionid"); e == nil {
			if cookie != nil {
				token = cookie.Value
			}
		}
	}
	if token == "" {
		return nil, false
	}
	log.Infoln("token:", token)

	sess, err := session.GetSessionById(token)
	if err != nil {
		log.Warningln("session ERR:", err.Error())
		return nil, false
	}
	log.Infoln("sess:", sess)

	if sess.Get("user") == nil {
		log.Errorln("session no user:", sess)
		return sess, false
	}

	return sess, true
}

func UserAdd(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	} else if r.Method != "POST" {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, "")
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
		gocommon.HttpErr(w, http.StatusInternalServerError, "用户中心错误.")
		return
	}

	log.Infoln("user add ok:", user.Userid)
	fmt.Fprintf(w, "{\"userid\":\"%v\"}", user.Userid)

	return
}

func UserModify(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	} else if r.Method != "POST" {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, "")
		return
	}

	//只有登录用户有权修改信息
	_, auth := AuthFilter(w, r)
	if auth == false {
		gocommon.HttpErr(w, http.StatusForbidden, "末登录用户.")
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
	} else if r.Method != "POST" {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, "")
		return
	}

	sess, err := session.GetSession(w, r, "")
	if err != nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, "会话错误.")
		log.Errorln("session.GetSession ERR:", err.Error())
		return
	}

	tmp := sess.Get("user")
	if tmp != nil {
		user := tmp.(*service.User)
		fmt.Fprintf(w, "{\"userid\":\"%v\", \"token\":\"%v\"}", user.Userid, sess.Id(""))
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
	log.Infoln("body:", string(body))

	user := &service.User{}
	if err := json.Unmarshal(body, user); err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		log.Errorln("json.Unmarshal(body, user) ERR: ", err, body)
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

	var has bool
	mUser := &service.User{}
	if user.Cellphone != "" {
		mUser.Cellphone = user.Cellphone
		has, err = mUser.LoginByCellphone()
	} else if user.Email != "" {
		mUser.Email = user.Email
		has, err = mUser.LoginByEmail()
	} else if user.Nickname != "" {
		mUser.Nickname = user.Nickname
		has, err = mUser.LoginByNickname()
	}

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

	loginPWD := service.EncryPWD(mUser.Userid, user.Password)
	if loginPWD != mUser.Password {
		gocommon.HttpErr(w, http.StatusForbidden, "用户密码不正确.")
		log.Warningln(*user, *mUser, "用户密码不正确.")
		return

	}

	sess.Set("user", mUser)
	log.Infoln("user login ok:", sess)
	fmt.Fprintf(w, "{\"userid\":\"%v\", \"token\":\"%v\"}", mUser.Userid, sess.Id(""))

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
	//	optionsFilter(w, r)
	//	if r.Method == "OPTIONS" {
	//		return
	//	}

	sess, auth := AuthFilter(w, r)
	if auth == false {
		gocommon.HttpErr(w, http.StatusForbidden, "末登录.")
		return
	}

	log.Infof("auth: %#v", sess)

	if mUser, ok := sess.Get("user").(*service.User); ok {
		mUser.Password = ""
		return
	}

	userStr, _ := json.Marshal(sess.Get("user"))
	w.Write(userStr)

	return
}
