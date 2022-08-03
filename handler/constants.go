package handler

import (
	"time"
)

type Action string

//Actions
const (
	JOIN_ACTION      Action = "join"
	SUBSCRIBE_ACTION Action = "subscribe"
	MESSAGE_ACTION   Action = "message"
	USERS_ACTION     Action = "users"
)

// Request
type Request struct {
	Action  Action                 `json:"action"`
	Payload map[string]interface{} `json:"payload"`
	ReqID   string                 `json:"reqId"`
}

//Response
type Response struct {
	Action  Action                 `json:"action"`
	Payload map[string]interface{} `json:"payload"`
	ReqID   string                 `json:"reqId"`
}

type Message struct {
	Recipient string    `json:"to"`
	Text      string    `json:"text"`
	Time      time.Time `json:"time"`
}
