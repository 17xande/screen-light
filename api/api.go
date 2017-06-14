package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type apiResponse struct {
	success bool
}

// ControlSend handles get requests with the rgb values
// sent in the querystring
func ControlSend(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	qs := r.URL.Query()
	green := "0"
	red := "0"
	blue := "0"

	if qs["r"] != nil {
		red = qs["r"][0]
	}
	if qs["g"] != nil {
		green = qs["g"][0]
	}
	if qs["b"] != nil {
		blue = qs["b"][0]
	}
	c := fmt.Sprintf("rgb(%s,%s,%s)", red, green, blue)
	// w.Write([]byte(c))
	hub.broadcast <- []byte(c)
	fmt.Println("API req:", c)

	w.WriteHeader(http.StatusOK)
	res := apiResponse{
		success: true,
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println("error encoding JSON: ", err)
	}
	return
}

// ColoursSave saves the colours in the controller interface
// to the colours.json file.
func ColoursSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	success := true

	f, err := os.Create("static/js/colours.json")
	defer f.Close()
	if err != nil {
		fmt.Printf("Could not open colours.json. %s\n", err)
		success = false
		err = nil
	}

	defer r.Body.Close()

	h, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	f.Write(h)
	// f.Sync()
	// fmt.Println(string(h))

	if success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	res := apiResponse{
		success: success,
	}

	err = json.NewEncoder(w).Encode(res)

	return
}
