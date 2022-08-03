package usermap

import "github.com/gorilla/websocket"

type UserMap struct {
	userMap map[string]*websocket.Conn
}

func New() *UserMap {
	return &UserMap{userMap: map[string]*websocket.Conn{}}
}

func (m *UserMap) IsUserPresent(userName string) bool {
	_, ok := m.userMap[userName]
	return ok
}

// AddUser adds or deletes user map
func (m *UserMap) AddUser(userChan chan *User) {
	for u := range userChan {
		m.userMap[u.Name] = u.Conn
	}
}

// DeleteUser deletes the user by connection
func (m *UserMap) DeleteUser(conn *websocket.Conn) {
	keyToBeDeleted := ""
	for k, v := range m.userMap {
		if v == conn {
			keyToBeDeleted = k
		}
	}
	delete(m.userMap, keyToBeDeleted)
}

// GetUserByConnection deletes the user by connection
func (m *UserMap) GetUserByConnection(conn *websocket.Conn) string {
	user := ""
	for k, v := range m.userMap {
		if v == conn {
			user = k
		}
	}
	return user
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
