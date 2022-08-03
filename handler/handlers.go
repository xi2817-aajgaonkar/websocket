package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xi2817-aajgaonkar/websocket/usermap"
)

func HandleUserAction(req *Request, u *usermap.UserMap, c *websocket.Conn) error {
	//get users
	users := u.GetUsers()
	out, err := json.Marshal(Response{
		Action: req.Action,
		ReqID:  req.ReqID,
		Payload: map[string]interface{}{
			"code":  200,
			"users": users,
		},
	})
	if err != nil {
		return err
	}
	if err = c.WriteMessage(websocket.TextMessage, out); err != nil {
		return err
	}
	return nil
}

func HandleMessageAction(req *Request, u *usermap.UserMap, c *websocket.Conn) error {
	msg := &Message{
		Recipient: req.Payload["to"].(string),
		Time:      time.Now(),
		Text:      req.Payload["text"].(string),
	}

	// check if user is present; of not then send corresponsing message to sender
	if !u.IsUserPresent(msg.Recipient) {
		out, err := json.Marshal(Response{
			Action: req.Action,
			ReqID:  req.ReqID,
			Payload: map[string]interface{}{
				"code":    400,
				"message": "recipient not present",
			},
		})
		if err != nil {
			return err
		}
		if err := c.WriteMessage(websocket.TextMessage, out); err != nil {
			return err
		}
		return errors.New("recipient not present")
	}

	// get conncection of recipient
	conn := u.GetConnection(msg.Recipient)
	if conn == nil {
		return errors.New("error: invalid recipent connection")
	}

	// send data to recipient
	jsonMsg, err := json.Marshal(Response{
		Action: SUBSCRIBE_ACTION,
		ReqID:  req.ReqID,
		Payload: map[string]interface{}{
			"code":    200,
			"message": msg.Text,
			"time":    msg.Time.String(),
			"from":    u.GetUserByConnection(c),
		},
	})
	if err != nil {
		return err
	}
	if err = conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
		return err
	}

	//send confirmation to sender
	out, err := json.Marshal(Response{
		Action: req.Action,
		ReqID:  req.ReqID,
		Payload: map[string]interface{}{
			"code":    200,
			"message": "message sent",
		},
	})
	if err != nil {
		return err
	}
	if err = c.WriteMessage(websocket.TextMessage, out); err != nil {
		return err
	}
	return nil
}

func HandleJoinAction(req *Request, userChannel chan *usermap.User, u *usermap.UserMap, c *websocket.Conn) error {
	users := u.GetUsers()
	fmt.Println(req)
	username := req.Payload["name"].(string)

	// check if user is not present already
	if u.IsUserPresent(username) {
		out, err := json.Marshal(Response{
			Action: req.Action,
			ReqID:  req.ReqID,
			Payload: map[string]interface{}{
				"code":    400,
				"message": "User already exists",
				"users":   users,
			},
		})
		if err != nil {
			return err
		}
		if err = c.WriteMessage(websocket.TextMessage, out); err != nil {
			return err
		}
		return errors.New("user Already present")
	}
	// append username to map for response
	users = append(users, username)

	// add user to map
	u.AddUser(&usermap.User{Name: username, Conn: c, Operation: usermap.ADD})
	out, err := json.Marshal(Response{
		Action: req.Action,
		ReqID:  req.ReqID,
		Payload: map[string]interface{}{
			"code":    200,
			"message": "new user added",
			"users":   users,
		},
	})
	if err != nil {
		return err
	}
	if err = c.WriteMessage(websocket.TextMessage, out); err != nil {
		return err
	}
	fmt.Println("Finished Join")
	return nil
}
