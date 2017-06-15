package api

import (
	"fmt"
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

// ServeWs serves the WebSocket connection
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
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

	// set connection limits
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

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
				fmt.Println(err)
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
