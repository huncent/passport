package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/controllers"
)

func init() {
	runtime.GOMAXPROCS(8)

	if e := common.InitDbPool("./db.conf"); e != nil {
		panic(e)
	}

	common.InitServ()
}

func main() {
	flag.Parse()

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/user/add", &controllers.UserAdd{})
	http.Handle("/user/login", &controllers.UserLogin{})
	http.Handle("/user/mod", &controllers.UserModify{})
	http.Handle("/user/auth", &controllers.UserAuth{})
	http.Handle("/user/logout", &controllers.UserLogout{})

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
