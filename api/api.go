package api

import (
	"fmt"
	"net/http"
)

// ControlSend handles get requests with the rgb values
// sent in the querystring
func ControlSend(hub *Hub, w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	c := fmt.Sprintf("rgb(%s,%s,%s)", qs["r"][0], qs["g"][0], qs["b"][0])
	// w.Write([]byte(c))
	hub.broadcast <- []byte(c)
	return
}

// ColoursSave saves the colours in the controller interface
// to the colours.json file.
func ColoursSave(w http.ResponseWriter, h *http.Request) {

}
