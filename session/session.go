package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var GlobalSessions *Manager

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    IProvider
	maxLifeTime int64
}

func NewManager(providerName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgetten import?", providerName)
	}
	return &Manager{
		provider:    provider,
		cookieName:  cookieName,
		maxLifeTime: maxLifeTime,
	}, nil
}

func (m *Manager) generateSessionId() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (m *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session ISession) {
	m.lock.Lock()
	defer m.lock.Unlock()
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		sid := m.generateSessionId()
		session, _ = m.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name:     m.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(m.maxLifeTime),
		}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = m.provider.SessionRead(sid)
	}
	return
}

func (m *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	var cookie *http.Cookie
	var err error
	cookie, err = r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.provider.SessionDestroy(cookie.Value)
	expiration := time.Now()
	cookie = &http.Cookie{
		Name:     m.cookieName,
		Path:     "/",
		HttpOnly: true,
		Expires:  expiration,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}

func (m *Manager) GC() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.provider.SessionGC(m.maxLifeTime)
	time.AfterFunc(time.Duration(m.maxLifeTime), func() {
		m.GC()
	})
}

type IProvider interface {
	SessionInit(sid string) (ISession, error)
	SessionRead(sid string) (ISession, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

var providers = make(map[string]IProvider)

func Register(name string, provider IProvider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	_, dup := providers[name]
	if dup {
		panic("session: Register called twice for provider" + name)
	}
	providers[name] = provider
}

type ISession interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionId() string
}
