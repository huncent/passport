package session_test

import (
	"testing"
	"time"

	"github.com/liuhengloveyou/passport/session"
)

var conf map[string]interface{} = map[string]interface{}{
	"store_type":    "mem",
	"cookie_name":   "passportsessionid",
	"idle_time":     86400,
	"cookie_expire": 86400,
}

func TestSession(t *testing.T) {

	session.InitDefaultSessionManager(conf)

	for {
		t.Log(">>>")
		time.Sleep(10)
		break
	}

	conf["cookie_name"] = "newName"
	session.InitDefaultSessionManager(conf)
}

func TestSessionKeys(t *testing.T) {
	session.InitDefaultSessionManager(conf)

	key := "demo"
	sess, _ := session.GetSessionById(&key)

	sess.Set("aa", "aaa")
	sess.Set("bb", "bbb")

	data := sess.Keys()
	t.Log("dump1:", data)

	sess.Set("ccc", "ccc")
	sess.Set("ddd", "ddd")
	sess.Delete("aa")

	t.Log("dump2:", data)
}
