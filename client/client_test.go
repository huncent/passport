package client_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/passport/client"
)

var passClient = &client.Passport{"http://localhost:8080"}

func TestUserAdd(t *testing.T) {
	data := `{"cellphone":"18510511015", "email":"liuhengloveyou@gmail.com", "nickname":"L", "password":"123456"}`
	stat, rst, e := passClient.UserAdd(data)
	fmt.Println(stat, rst, e)
}

func TestUserLogin(t *testing.T) {
	data := `{"cellphone":"18510511015", "email":"liuhengloveyou@gmail.com", "nickname":"L", "password":"123456"}`
	stat, rst, e := passClient.UserLogin(data)
	fmt.Println(stat, rst, e)
}
