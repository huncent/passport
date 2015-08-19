package session_test

import (
	"testing"
	"time"

	"github.com/liuhengloveyou/passport/session"
)

func TestSession(t *testing.T) {
	conf := map[string]interface{}{
		"store_type":    "mem",
		"cookie_name":   "passportsessionid",
		"idle_time":     259200,
		"cookie_expire": 259200,
	}

	session.InitDefaultSessionManager(conf)

	for {
		t.Log(">>>")
		time.Sleep(10)
		break
	}

	conf["cookie_name"] = "newName"
	session.InitDefaultSessionManager(conf)

	for {
		t.Log(">>>")
		time.Sleep(10)
	}
}
