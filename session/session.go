package session

import (
	"container/list"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	log "github.com/golang/glog"
)

var defaultSessionManager *SessionManager = nil

var stores = make(map[string]SessionStoreType)

type SessionStoreType func(interface{}) (SessionStore, error)

type SessionStore interface {
	Id(string) string
	Active(set bool) int64
	Get(key interface{}) interface{}
	Set(key, val interface{}) error
	Keys() []interface{}
	Delete(key interface{}) error
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
	StoreType    string      `json:"store_type"`
	CookieName   string      `json:"cookie_name"`
	IdleTime     int64       `json:"idle_time"`
	CookieExpire int         `json:"cookie_expire"`
	Domain       string      `json:"domain"`
	StoreConfig  interface{} `json:"store_config"`

	sessions  map[string]*list.Element
	list      *list.List
	lock      sync.RWMutex
	destroied bool
}

func NewSessionManager(sessionConfig interface{}) (m *SessionManager) {
	if sessionConfig == nil {
		return nil
	}

	m = &SessionManager{}

	var byteConf []byte
	var err error
	if byteConf, err = json.Marshal(sessionConfig); err != nil {
		return nil
	}

	if err = json.Unmarshal(byteConf, m); err != nil {
		return nil
	}

	m.sessions = make(map[string]*list.Element)
	m.list = list.New()
	m.gc()

	return m
}

func (p *SessionManager) GetSessionById(sessionid *string) (session SessionStore, err error) {
	sid := ""

	if sessionid == nil {
		if sid, err = p.sessionId(); err != nil {
			return nil, err
		}
	} else {
		sid = *sessionid
	}

	if sess, ok := p.sessions[sid]; ok {
		session = sess.Value.(SessionStore)
		p.lock.Lock()
		p.list.MoveToBack(sess)
		p.lock.Unlock()
		return
	}

	// 新会话
	session, err = newSessionStore(p.StoreType, p.StoreConfig)
	if err != nil {
		return
	}
	session.Id(sid)

	p.lock.Lock()
	p.sessions[sid] = p.list.PushBack(session)
	p.lock.Unlock()

	return
}

func (p *SessionManager) GetSession(w http.ResponseWriter, r *http.Request, sessionid *string) (session SessionStore, err error) {
	writeCookie := false
	sid := ""

	cookie, errs := r.Cookie(p.CookieName)
	if errs != nil || cookie.Value == "" {
		if sessionid == nil {
			sid, err = p.sessionId()
		} else {
			sid = *sessionid
		}
		writeCookie = true
	} else {
		sid, err = url.QueryUnescape(cookie.Value)
	}
	if err != nil {
		return
	}

	if sess, ok := p.sessions[sid]; ok {
		session = sess.Value.(SessionStore)
		p.lock.Lock()
		p.list.MoveToBack(sess)
		p.lock.Unlock()
		return
	}

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
			Name:   p.CookieName,
			Value:  url.QueryEscape(sid),
			Path:   "/",
			Domain: p.Domain,
		}

		if p.CookieExpire >= 0 {
			cookie.MaxAge = p.CookieExpire
		}

		http.SetCookie(w, cookie)
	}

	r.AddCookie(cookie)

	return
}

func (p *SessionManager) Destroy() {
	p.sessions = nil
	p.list = nil
	p.destroied = true
}

func (p *SessionManager) SessionDestroy(w http.ResponseWriter, r *http.Request) (userid int64, sessionid string) {
	cookie, err := r.Cookie(p.CookieName)
	if err != nil || cookie.Value == "" {
		return
	}

	sessionid, _ = url.QueryUnescape(cookie.Value)

	if session, ok := p.sessions[sessionid]; ok {
		if session.Value.(SessionStore).Get("id") == nil {
			return
		}

		userid = session.Value.(SessionStore).Get("id").(int64)
		session.Value.(SessionStore).Release()
	}

	http.SetCookie(w, &http.Cookie{
		Name:     p.CookieName,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(),
		MaxAge:   -1})

	return
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
	var sleep int64 = 10

	for p.destroied == false {
		var element *list.Element

		p.lock.RLock()
		if element = p.list.Front(); element == nil {
			p.lock.RUnlock()
			break
		}

		if (element.Value.(SessionStore).Active(false) + p.IdleTime) > time.Now().Unix() {
			sleep = (element.Value.(SessionStore).Active(false) + int64(p.IdleTime)) - time.Now().Unix()
			p.lock.RUnlock()
			break
		}
		p.lock.RUnlock()

		log.Warningln("sessionrelease:", element.Value.(SessionStore).Id(""))
		p.lock.Lock()
		element.Value.(SessionStore).Release()
		delete(p.sessions, element.Value.(SessionStore).Id(""))
		p.list.Remove(element)
		p.lock.Unlock()
	}

	if p.destroied == false {
		time.AfterFunc(time.Duration(sleep)*time.Second, p.gc)
	}
}

// 公开接口
func InitDefaultSessionManager(conf interface{}) *SessionManager {
	if defaultSessionManager != nil {
		defaultSessionManager.Destroy()
	}

	defaultSessionManager = NewSessionManager(conf)
	return defaultSessionManager
}

func GetSession(w http.ResponseWriter, r *http.Request, sessionid *string) (session SessionStore, err error) {
	return defaultSessionManager.GetSession(w, r, sessionid)
}

func GetSessionById(sessionid *string) (session SessionStore, err error) {
	return defaultSessionManager.GetSessionById(sessionid)
}

func SessionDestroy(w http.ResponseWriter, r *http.Request) (userid int64, sessionid string) {
	return defaultSessionManager.SessionDestroy(w, r)
}
