// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "3.7.100.88:7000", "http service address")
var interrupt = make(chan os.Signal, 1)

func main() {
	flag.Parse()
	log.SetFlags(0)

	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
		}
	}()

	payload := map[string]interface{}{"name": "john"}
	final := map[string]interface{}{"action": "join", "payload": payload, "reqId": "thisisreqID123"}
	// final := map[string]interface{}{"action": "join", "payload": payload, "reqId": "thisisreqID123"}
	out, err := json.Marshal(final)
	if err != nil {
		panic(err)
	}
	log.Println("Writing socket")
	err = c.WriteMessage(websocket.TextMessage, []byte(out))
	if err != nil {
		log.Println("write:", err)
		return
	}

	payload = map[string]interface{}{"message": "hii atharva", "to": "atharva"}
	final = map[string]interface{}{"action": "message", "payload": payload, "reqId": "thisisreqID123"}
	// final := map[string]interface{}{"action": "join", "payload": payload, "reqId": "thisisreqID123"}
	out, err = json.Marshal(final)
	if err != nil {
		panic(err)
	}
	log.Println("Writing socket")
	err = c.WriteMessage(websocket.TextMessage, []byte(out))
	if err != nil {
		log.Println("write:", err)
		return
	}

	ticker := time.NewTicker(time.Second * 50000)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// payload := map[string]interface{}{"name": "atharva"}
			payload := map[string]interface{}{"name": "atharva"}
			final := map[string]interface{}{"action": "join", "payload": payload, "reqId": "thisisreqID123"}
			// final := map[string]interface{}{"action": "join", "payload": payload, "reqId": "thisisreqID123"}
			out, err := json.Marshal(final)
			if err != nil {
				panic(err)
			}
			log.Println("Writing socket")
			err = c.WriteMessage(websocket.TextMessage, []byte(out))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
