package session_test

import (
	"testing"
	"time"

	"github.com/liuhengloveyou/passport/session"
)

var conf map[string]interface{} = map[string]interface{}{
	"store_type":    "mem",
	"cookie_name":   "passportsessionid",
	"idle_time":     10,
	"cookie_expire": 10,
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
	t.Log("keys1:", data)

	sess.Set("ccc", "ccc")
	sess.Set("ddd", "ddd")
	sess.Delete("aa")

	t.Log("keys2:", data)
}

func TestSessionGC(t *testing.T) {

	pf := func(sess session.SessionStore) {
		t.Log("PrepireRelease:", sess)
	}

	session.InitDefaultSessionManager(conf)
	session.SetPrepireRelease(pf)
	key := "demo"
	sess, _ := session.GetSessionById(&key)

	sess.Set("aa", "aaa")
	sess.Set("bb", "bbb")

	time.Sleep(6 * time.Second)
	t.Log("aa:", sess.Get("aa"))
	t.Log("bb:", sess.Get("bb"))

	time.Sleep(6 * time.Second)
	t.Log("aa:", sess.Get("aa"))
	t.Log("bb:", sess.Get("bb"))
}
