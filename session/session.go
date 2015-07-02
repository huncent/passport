package session

import (
	"container/list"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var defaultSessionManager *SessionManager = NewSessionManager(nil)

var stores = make(map[string]SessionStoreType)

type SessionStoreType func(interface{}) (SessionStore, error)

type SessionStore interface {
	Id(string) string
	Get(key interface{}) interface{}
	Set(key, val interface{}) error
	Delete(key interface{}) error
	Active() int64
	Release()
}

func RegisterSessionStore(name string, one SessionStoreType) {
	if one == nil {
		panic("Register SessionStore nil")
	}

	if _, dup := stores[name]; dup {
		panic("Register SessionStore duplicate for " + name)
	}

	stores[name] = one
}

func newSessionStore(typeName string, config interface{}) (SessionStore, error) {
	if newFunc, ok := stores[typeName]; ok {
		return newFunc(config)
	}

	return nil, fmt.Errorf("No SessionManager types " + typeName)
}

////
type SessionManager struct {
	StoreType    string `json:"store_type"`
	CookieName   string `json:"cookie_name"`
	IdleTime     int64  `json:"idle_time"`
	CookieExpire int    `json:"cookie_expire"`
	Domain       string `json:"domain"`
	StoreConfig  string `json:"store_config"`

	sessions map[string]*list.Element
	list     *list.List
	lock     sync.RWMutex
}

func NewSessionManager(sessionConfig interface{}) (m *SessionManager) {
	m = &SessionManager{
		StoreType:    "mem",
		CookieName:   "gosessionid",
		IdleTime:     10,
		CookieExpire: 0,
		sessions:     make(map[string]*list.Element),
		list:         list.New(),
	}

	m.gc()

	return
}

func (p *SessionManager) GetSession(w http.ResponseWriter, r *http.Request) (session SessionStore, err error) {
	sid := ""
	writeCookie := false

	cookie, errs := r.Cookie(p.CookieName)
	if errs != nil || cookie.Value == "" {
		sid, err = p.sessionId()
		writeCookie = true
	} else {
		sid, err = url.QueryUnescape(cookie.Value)
	}

	if err != nil {
		return
	}

	p.lock.RLock()
	if sess, ok := p.sessions[sid]; ok {
		session = sess.Value.(SessionStore)

		// @@@ 不好
		p.lock.RUnlock()
		p.lock.Lock()
		p.list.MoveToBack(sess)
		p.lock.Lock()
		p.lock.RLock()

		return
	}
	p.lock.RUnlock()

	// 新会话
	session, err = newSessionStore(p.StoreType, p.StoreConfig)
	if err != nil {
		return
	}
	session.Id(sid)

	p.lock.Lock()
	p.sessions[sid] = p.list.PushBack(session)
	p.lock.Unlock()

	if writeCookie == true {
		cookie = &http.Cookie{
			Name:     p.CookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			Domain:   p.Domain,
		}

		if p.CookieExpire >= 0 {
			cookie.MaxAge = p.CookieExpire
		}

		http.SetCookie(w, cookie)
	}

	r.AddCookie(cookie)

	return
}

func (p *SessionManager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(p.CookieName)
	if err != nil || cookie.Value == "" {
		return
	}

	sid, _ := url.QueryUnescape(cookie.Value)
	p.lock.Lock()
	if session, ok := p.sessions[sid]; ok {
		session.Value.(SessionStore).Release()
		delete(p.sessions, sid)
		p.list.Remove(session)
	}
	p.lock.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     p.CookieName,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(),
		MaxAge:   -1})
}

func (p *SessionManager) sessionId() (string, error) {
	b := make([]byte, 24)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", fmt.Errorf("Could not successfully read from the system CSPRNG.")
	}

	return hex.EncodeToString(b), nil
}

func (p *SessionManager) gc() {
	var sleep int64 = 3

	for {
		var element *list.Element

		p.lock.RLock()
		if element = p.list.Front(); element == nil {
			p.lock.RUnlock()
			break
		}

		if (element.Value.(SessionStore).Active() + p.IdleTime) > time.Now().Unix() {
			sleep = (element.Value.(SessionStore).Active() + int64(p.IdleTime)) - time.Now().Unix()
			p.lock.RUnlock()
			break
		}
		p.lock.RUnlock()

		p.lock.Lock()
		delete(p.sessions, element.Value.(SessionStore).Id(""))
		p.list.Remove(element)
		p.lock.Unlock()
	}

	time.AfterFunc(time.Duration(sleep)*time.Second, p.gc)
}

// 公开接口
func GetSession(w http.ResponseWriter, r *http.Request) (session SessionStore, err error) {
	return defaultSessionManager.GetSession(w, r)
}
