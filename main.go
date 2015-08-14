package main

import (
	"flag"
	"runtime"

	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/face"
)

func init() {
	runtime.GOMAXPROCS(8)

	if e := common.InitPassportServ(); e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse()

	face.HttpService()
}
