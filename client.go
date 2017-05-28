package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message
	writeWait = 1 * time.Second

	// Time allowed to read the next pong message from the web client
	pongWait = 30 * time.Second

	// Send pings to web client with this period. Must be less than pongWait
	pingPeriod = pongWait * 9 / 10

	// Max message size allowed from web client
	maxMessageSize = 128
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  256,
	WriteBufferSize: 256,
}

// Client is a middleman between the websocket connection and the hub.
// It represents each connected web client
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	// Buffered channel of outbound messages
	send chan []byte
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}

// readPump pumps messages from the websocket connection to the hub
// This is run through a goroutine for each connection.
func (c *Client) readPump() {
	// close this client's connection and remove them from the map once this
	// function exits. It runs an infinite look that only breaks when the
	// web client unregisters.
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// set connection limits
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// not sure why we have to set the read deadline everytime we receive a pong?
	// Perhaps because there is no timeout, just a deadline that needs to keep being
	// updated everytime a new pong is received to keep the connection alive?
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// infinite loop listening for messages from the web client
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			// break out of the loop if there is an error reading the message
			break
		}
		c.hub.broadcast <- message
	}
}

// writePump sends the messages from the hub to the websocket connections.
// a goroutine running writePump is started for each connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// the hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// gets the next available writer
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// add queued chat messages to the current websocket message
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}
		// send keep-alive ping signal to the websocket every time
		// the ticker sends a message on the channel
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}
}

// TODO:
// 1. change to send bytes instead of text
// 2. change buffers to only handle the length of colour values, eg #111111
// 3. See what causes multiple values to be sent along on one message/ data thing.
// if that can't be avoided, discard early values and only process latest value.
