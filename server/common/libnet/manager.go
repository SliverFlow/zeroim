package libnet

import (
	"github.com/SliverFlow/zeroim/server/common/hash"
	"github.com/SliverFlow/zeroim/server/common/session"
	"sync"
)

const sessionMapNum = 32

// Manager is a socket manager.
type Manager struct {
	Name        string                    // 服务器名称
	sessionMaps [sessionMapNum]sessionMap //
	disposeFlag bool
	disposeOnce sync.Once
	disposeWait sync.WaitGroup
}

// sessionMap is a session map.
type sessionMap struct {
	sync.RWMutex
	sessions      map[session.Session]*Session
	tokenSessions map[string][]session.Session
}

// NewManager
func NewManager(name string) *Manager {
	manager := &Manager{
		Name: name,
	}

	for i := 0; i < sessionMapNum; i++ {
		manager.sessionMaps[i].sessions = make(map[session.Session]*Session)
		manager.sessionMaps[i].tokenSessions = make(map[string][]session.Session)
	}

	return manager
}

// GetSession
func (m *Manager) GetSession(sessionId session.Session) *Session {
	token := sessionId.Token()
	hashId := hash.Hash([]byte(token))
	smap := &m.sessionMaps[hashId&sessionMapNum]
	smap.RLock()

	defer smap.RUnlock()

	return smap.sessions[sessionId]
}

// GetTokenSessions
func (m *Manager) GetTokenSessions(token string) []*Session {
	hashId := hash.Hash([]byte(token))
	smap := &m.sessionMaps[hashId&sessionMapNum]
	smap.RLock()
	defer smap.RUnlock()

	sessionIds := smap.tokenSessions[token]

	var sessions []*Session
	for _, sessionId := range sessionIds {
		sessions = append(sessions, smap.sessions[sessionId])
	}
	return sessions
}

func (m *Manager) AddSession(session *Session) {
	sessionId := session.Session()
	token := session.token
	hashId := hash.Hash([]byte(token))
	smap := &m.sessionMaps[hashId%sessionMapNum]

	smap.Lock()
	defer smap.Unlock()

	smap.sessions[sessionId] = session
	smap.tokenSessions[token] = append(smap.tokenSessions[token], sessionId)
}

func (m *Manager) removeSession(session *Session) {

}

func (m *Manager) Close() {
	m.disposeOnce.Do(func() {
		m.disposeFlag = true

		for i := 0; i < sessionMapNum; i++ {
			smap := &m.sessionMaps[i]
			smap.Lock()
			for _, session := range smap.sessions {
				session.Close()
			}
			smap.Unlock()
		}
		m.disposeWait.Wait()
	})
}
