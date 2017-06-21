package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"fmt"

	"github.com/gorilla/websocket"
)

// Controller is a middleman between the websocket connection to the controller
// interface and the hub.
type Controller struct {
	hub  *Hub
	conn *websocket.Conn
	// channel of outbound instructions
	instruct chan []byte
	// channel of inbound connection messages from hub
	conns chan int
}

// ServeController serves the controller interface
func ServeController(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	co := &Controller{
		hub:      hub,
		conn:     conn,
		instruct: make(chan []byte, 256),
		conns:    make(chan int, 0),
	}
	co.hub.regController <- co

	// set connection limits
	co.conn.SetReadLimit(maxMessageSize)
	co.conn.SetReadDeadline(time.Now().Add(pongWait))
	co.conn.SetPongHandler(func(string) error {
		co.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go co.socketWrite()
	co.socketRead()
}

// ControlPreset handles requests sent to engage a preset
// Used in the REST API
func ControlPreset(hub *Hub, w http.ResponseWriter, r *http.Request) {

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

func (co *Controller) socketWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		co.conn.Close()
	}()

	for {
		select {
		case conns, ok := <-co.conns:
			co.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				co.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := co.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println(err)
				return
			}

			w.Write([]byte(strconv.Itoa(conns)))
			if err := w.Close(); err != nil {
				fmt.Println(err)
				return
			}
		case <-ticker.C:
			co.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := co.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
