package main

// Hub maintains the map of active clients and broadcasts messages to the clients.
type Hub struct {
	clients map[*Client]bool
	// register requests from clients
	register chan *Client
	// unregister requests from clients
	unregister chan *Client
	// inbound messages from clients
	broadcast chan []byte
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) run() {
	for {
		// select prevents blocking by creating a way to drop messages
		// from channels that are full. It looks like a switch statement
		select {
		// handle a message from the Hub's register channel, store the message
		// which in this case is a *Client, in the client variable
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			// check if the client is in the list. If it's not, ignore this request.
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// loop through all clients and send the message to them all
			for client := range h.clients {
				// send the message to the clients. Use a select in case the client
				// has disconnected
				select {
				// if this channel is unavailable or full we assume the connection
				// has closed so we just remove the client.
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
