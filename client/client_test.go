package client_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/passport/client"
)

var passClient = &client.Passport{"http://localhost:8080"}

func TestUserAdd(t *testing.T) {
	stat, rst, e := passClient.UserAdd("18510511015", "liuhengloveyou@gmail.com", "123456")
	fmt.Println(stat, rst, e)
}

func TestUserLogin(t *testing.T) {
	stat, rst, e := passClient.UserLogin("", "18510511015", "liuhengloveyou@gmail.com", "123456")
	fmt.Println(stat, rst, e)
}
