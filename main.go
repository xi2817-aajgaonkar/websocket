// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"os/signal"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var interrupt = make(chan os.Signal, 1)

type Message struct {
	Operation string `json:"operation"`
	Message   string `json:"message"`
}

func subscribe(w http.ResponseWriter, r *http.Request, c *websocket.Conn) {
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			out, err := json.Marshal(Message{Message: t.String(), Operation: "subscribe"})
			if err != nil {
				log.Println("write:", err)
				return
			}
			err = c.WriteMessage(websocket.TextMessage, out)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func echo(w http.ResponseWriter, r *http.Request, c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	temp := &Message{}
	err = c.ReadJSON(temp)
	if err != nil {
		log.Println("read:", err)
	}

	if temp.Operation == "echo" {
		echo(w, r, c)
	}

	if temp.Operation == "subscribe" {
		subscribe(w, r, c)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
