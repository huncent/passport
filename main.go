package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"

	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/face"
)

var (
	confile = flag.String("c", "passport.conf.sample", "配置文件路径.")

	initSys = flag.Bool("init", false, "初始化系统.")
)

func main() {
	flag.Parse()

	if e := common.InitPassportServ(*confile); e != nil {
		panic(e)
	}

	if *initSys {
		if e := common.InitSystem(common.ServConfig.DBs); e != nil {
			panic(e)
		}
		return
	}

	switch common.ServConfig.Face {
	case "http":
		face.HttpService()
	case "grpc":
		face.GrpcFace()
	default:
		fmt.Println("face: [http | grpc]")
	}
}
