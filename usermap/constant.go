package usermap

import "github.com/gorilla/websocket"

type User struct {
	Name      string
	Conn      *websocket.Conn
	Operation Operation
}

type Operation string

// Operation
const (
	ADD    Operation = "add"
	DELETE Operation = "delete"
)
