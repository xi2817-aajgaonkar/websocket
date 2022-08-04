package usermap

import "github.com/gorilla/websocket"

type User struct {
	Name string
	Conn *websocket.Conn
}
