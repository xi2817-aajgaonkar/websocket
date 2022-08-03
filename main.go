// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/xi2817-aajgaonkar/websocket/handler"
	"github.com/xi2817-aajgaonkar/websocket/usermap"
)

// use default options
var addr = flag.String("addr", ":7000", "http service address")

func wsHandlers(u *usermap.UserMap, userChannel chan *usermap.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Start handlers")
		// upgrade connection to websocket
		var upgrader = websocket.Upgrader{}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		// close connection and delete chanel
		defer func() {
			fmt.Println("Deleting User", u.GetUsers())
			c.Close()
			u.DeleteUser(c)
		}()

		for {
			req := &handler.Request{}
			err = c.ReadJSON(req)
			if err != nil {
				log.Println("read:", err)
			}

			switch req.Action {
			case handler.JOIN_ACTION:
				fmt.Println("inside join")
				if err := handler.HandleJoinAction(req, userChannel, u, c); err != nil {
					log.Println("Error in JOIN ACTION", err)
				}
			//write users in response to this request
			case handler.USERS_ACTION:
				fmt.Println("inside user")
				if err := handler.HandleUserAction(req, u, c); err != nil {
					log.Println("Error in USER ACTION", err)
				}

			case handler.MESSAGE_ACTION:
				fmt.Println("inside message")
				// get user from map and send data to that connection
				if err := handler.HandleMessageAction(req, u, c); err != nil {
					log.Println("Error in Message ACTION ", err)
				}
			}
		}
	}
}

func main() {

	var userChannel = make(chan *usermap.User)
	u := usermap.New()
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/ws", wsHandlers(u, userChannel))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
