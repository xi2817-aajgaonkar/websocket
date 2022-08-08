// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/xi2817-aajgaonkar/websocket/handler"
	"github.com/xi2817-aajgaonkar/websocket/usermap"
)

// use default options
var addr = flag.String("addr", ":7000", "http service address")

func sendUsersToAllConnections(u *usermap.UserMap) error {
	out, err := json.Marshal(handler.Response{
		Action: handler.USERS_ACTION,
		ReqID:  "downsream-request-id",
		Payload: map[string]interface{}{
			"code":    200,
			"message": "new user added",
			"users":   u.GetUsers(),
		},
	})
	if err != nil {
		return err
	}

	for _, r := range u.GetUsers() {
		conn := u.GetConnection(r)
		if err = conn.WriteMessage(websocket.TextMessage, out); err != nil {
			return err
		}
	}

	return nil
}

func HandleWsHandlers(u *usermap.UserMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		userName := ""
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		// close connection and delete chanel
		defer func() {
			c.Close()
			u.DeleteUser(userName)
			if err := sendUsersToAllConnections(u); err != nil {
				log.Print("write:", err)
				return
			}
		}()

		for {
			req := &handler.Request{}
			err = c.ReadJSON(req)
			if err != nil {
				log.Println("read:", err)
			}

			switch req.Action {
			case handler.JOIN_ACTION:
				if err := handler.HandleJoinAction(req, &userName, u, c); err != nil {
					log.Println("Error in JOIN ACTION", err)
					return
				}
				if err := sendUsersToAllConnections(u); err != nil {
					log.Print("write:", err)
					return
				}
			//write users in response to this request
			case handler.USERS_ACTION:
				if err := handler.HandleUserAction(req, u, c); err != nil {
					log.Println("Error in USER ACTION", err)
				}

			case handler.MESSAGE_ACTION:
				// get user from map and send data to that connection
				if err := handler.HandleMessageAction(req, userName, u, c); err != nil {
					log.Println("Error in Message ACTION ", err)
				}
			}
		}
	}
}

// CreateCorsObject creates a cors object with the required config
func createCorsObject() *cors.Cors {
	return cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc: func(s string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		ExposedHeaders: []string{"Authorization", "Content-Type"},
	})
}

func main() {
	u := usermap.New()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ws", HandleWsHandlers(u))
	flag.Parse()
	log.SetFlags(0)

	corsObj := createCorsObject()
	Handler := corsObj.Handler(router)
	http.ListenAndServe(":7000", Handler)
}
