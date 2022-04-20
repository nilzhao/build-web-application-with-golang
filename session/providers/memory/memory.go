package memory

import (
	"container/list"
	"sync"
	"time"

	"github.com/nilzhao/build-web-application-with-golang/session"
)

var pder = &Provider{list: list.New()}

type SessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[any]any
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) Get(key any) any {
	pder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	}
	return nil
}

func (st *SessionStore) Delete(key any) error {
	delete(st.value, key)
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionId() string {
	return st.sid
}

type Provider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func (p *Provider) SessionInit(sid string) (session.ISession, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[any]any, 0)
	newSessionStore := &SessionStore{
		sid:          sid,
		timeAccessed: time.Now(),
		value:        v,
	}
	element := pder.list.PushFront(newSessionStore)
	pder.sessions[sid] = element
	return newSessionStore, nil
}

func (p *Provider) SessionRead(sid string) (session.ISession, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	}
	sess, err := pder.SessionInit(sid)
	return sess, err
}

func (p *Provider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		p.list.Remove(element)
		return nil
	}
	return nil
}

func (p *Provider) SessionGC(maxLifeTime int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for {
		element := p.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
			p.list.Remove(element)
			delete(p.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (p *Provider) SessionUpdate(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if element, ok := p.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		p.list.MoveToFront(element)
		return nil

	}
	return nil
}

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", pder)
}
