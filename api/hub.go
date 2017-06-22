package api

// Hub maintains the map of active clients and broadcasts messages to the clients.
type Hub struct {
	clients         map[*Client]bool
	register        chan *Client // register requests from clients
	unregister      chan *Client // unregister requests from clients
	controllers     map[*Controller]bool
	regController   chan *Controller
	unregController chan *Controller
	broadcast       chan []byte // inbound messages from clients
	lastInstruction []byte
}

// NewHub returns a new *Hub
func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		controllers:     make(map[*Controller]bool),
		regController:   make(chan *Controller),
		unregController: make(chan *Controller),
		broadcast:       make(chan []byte),
	}
}

// Run starts the listening loop of the Hub
func (h *Hub) Run() {
	for {
		// select prevents blocking by creating a way to drop messages
		// from channels that are full. It looks like a switch statement
		select {
		// handle a message from the Hub's register channel, in this case it's
		// a new *Client. Store it in the hub.
		case client := <-h.register:
			h.clients[client] = true
			notifyControllers(h)
			if h.lastInstruction != nil {
				client.send <- h.lastInstruction
			}
		case client := <-h.unregister:
			// check if the client is in the list. If it's not, ignore this request.
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				notifyControllers(h)
			}
		case controller := <-h.regController:
			h.controllers[controller] = true
		case controller := <-h.unregController:
			if _, ok := h.controllers[controller]; ok {
				delete(h.controllers, controller)
				close(controller.instruct)
			}
		case message := <-h.broadcast:
			// store the last instruction to send to new connections
			h.lastInstruction = message
			// loop through all clients and send the message to them all
			for client := range h.clients {
				// send the message to the clients. Use a select in case the client
				// has disconnected
				select {
				case client.send <- message:
				// if this channel is unavailable or full we assume the connection
				// has closed so we just remove the client.
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func notifyControllers(h *Hub) {
	for controller := range h.controllers {
		select {
		case controller.conns <- len(h.clients):

		default:
			close(controller.conns)
			delete(h.controllers, controller)
		}
	}
}
