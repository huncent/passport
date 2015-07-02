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
	p.active = time.Now().Unix()

	p.lock.Lock()
	p.data[key] = val
	p.lock.Unlock()

	return nil
}

func (p *MemSessionStore) Get(key interface{}) interface{} {
	p.active = time.Now().Unix()

	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.data[key]

}

func (p *MemSessionStore) Delete(key interface{}) error {
	p.active = time.Now().Unix()

	p.lock.Lock()
	delete(p.data, key)
	p.lock.Unlock()

	return nil
}

func (p *MemSessionStore) Active() int64 {
	return p.active
}

func (p *MemSessionStore) Release() {}

func NewMemSessionStore(config interface{}) (SessionStore, error) {
	return &MemSessionStore{
		data:   make(map[interface{}]interface{}),
		active: time.Now().Unix()}, nil
}

func init() {
	RegisterSessionStore("mem", NewMemSessionStore)
}
