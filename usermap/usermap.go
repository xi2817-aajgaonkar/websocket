package usermap

import (
	"sync"

	"github.com/gorilla/websocket"
)

type UserMap struct {
	userMap map[string]*websocket.Conn
	mutex   sync.RWMutex
}

func New() *UserMap {
	return &UserMap{userMap: map[string]*websocket.Conn{}, mutex: sync.RWMutex{}}
}

func (m *UserMap) IsUserPresent(userName string) bool {
	_, ok := m.userMap[userName]
	return ok
}

// AddUser adds or deletes user map
func (m *UserMap) AddUser(u *User) {
	m.mutex.Lock()
	m.userMap[u.Name] = u.Conn
	m.mutex.Unlock()
}

// DeleteUser deletes the user by connection
func (m *UserMap) DeleteUser(userName string) {
	m.mutex.Lock()
	delete(m.userMap, userName)
	m.mutex.Unlock()
}

//GetUsers returns active users
func (m *UserMap) GetUsers() []string {
	keys := make([]string, 0, len(m.userMap))
	for k := range m.userMap {
		keys = append(keys, k)
	}
	return keys
}

//GetUsers returns active users
func (m *UserMap) GetConnection(userName string) *websocket.Conn {
	if conn, ok := m.userMap[userName]; ok {
		return conn
	}
	return nil
}
