package session

import (
	"sync"
	"time"
)

type MemSessionStore struct {
	id     string
	data   map[interface{}]interface{}
	active int64
	lock   sync.RWMutex
}

func (p *MemSessionStore) Id(id string) string {
	if id != "" {
		p.id = id
	}

	return p.id
}

func (p *MemSessionStore) Set(key, val interface{}) error {
	p.lock.Lock()
	p.data[key] = val
	p.lock.Unlock()

	return nil
}

func (p *MemSessionStore) Get(key interface{}) (v interface{}) {
	p.lock.RLock()
	v = p.data[key]
	p.lock.RUnlock()

	return
}

func (p *MemSessionStore) Delete(key interface{}) error {
	p.lock.Lock()
	delete(p.data, key)
	p.lock.Unlock()

	return nil
}

func (p *MemSessionStore) Keys() []interface{} {
	len := len(p.data)
	tmp := make([]interface{}, len)

	p.lock.RLock()
	for k, _ := range p.data {
		len--
		tmp[len] = k
	}
	p.lock.RUnlock()

	return tmp
}

func (p *MemSessionStore) Active(set bool) (val int64) {
	val = p.active

	if set {
		p.lock.Lock()
		p.active = time.Now().Unix()
		p.lock.Unlock()
	}

	return
}

func (p *MemSessionStore) Release() {
	p.lock.Lock()
	p.active = -1
	p.data = nil
	p.lock.Unlock()
}

func NewMemSessionStore(config interface{}) (SessionStore, error) {
	return &MemSessionStore{
		data:   make(map[interface{}]interface{}),
		active: time.Now().Unix()}, nil
}

func init() {
	RegisterSessionStore("mem", NewMemSessionStore)
}
