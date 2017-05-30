package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Controller is a middleman between the websocket connection to the controller
// interface and the hub.
type Controller struct {
	hub  *Hub
	conn *websocket.Conn
	// channel of outbound instructions
	instruct chan []byte
	// channel of inbound status messages from hub
	messages chan []byte
}

// serve the controller interface
func serveController(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	co := &Controller{
		hub:      hub,
		conn:     conn,
		instruct: make(chan []byte, 256),
	}
	co.hub.regController <- co
	// set connection limits
	co.conn.SetReadLimit(maxMessageSize)
	co.conn.SetReadDeadline(time.Now().Add(pongWait))
	co.conn.SetPongHandler(func(string) error {
		co.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

func (co *Controller) socketRead() {
	defer func() {
		co.hub.unregController <- co
		co.conn.Close()
	}()

	for {
		_, instruction, err := co.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		co.hub.broadcast <- instruction
	}
}
