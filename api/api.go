package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type apiResponse struct {
	Success   bool   `json:"success"`
	Preset    int    `json:"preset,omitempty"`
	Color     string `json:"color"`
	Animation string `json:"animation,omitempty"`
	Frequency int    `json:"frequency,omitempty"`
}

// ControlSend handles get requests with the rgb values
// sent in the querystring
func ControlSend(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	qs := r.URL.Query()
	green := "0"
	red := "0"
	blue := "0"

	ar := apiResponse{
		Success: true,
	}

	if qs["r"] != nil {
		red = qs["r"][0]
	}
	if qs["g"] != nil {
		green = qs["g"][0]
	}
	if qs["b"] != nil {
		blue = qs["b"][0]
	}
	if qs["a"] != nil {
		ar.Animation = qs["a"][0]
	}
	if qs["f"] != nil && len(qs["f"][0]) > 0 {
		f, err := strconv.Atoi(qs["f"][0])
		ar.Frequency = f
		if err != nil {
			ar.Success = false
			log.Printf("error: %v", err)
		}
	}
	if qs["p"] != nil && len(qs["p"][0]) > 0 {
		p, err := strconv.Atoi(qs["p"][0])
		ar.Preset = p
		if err != nil {
			ar.Success = false
			log.Printf("error: %v", err)
		}
	}

	ar.Color = fmt.Sprintf("rgb(%s,%s,%s)", red, green, blue)
	jsAr, err := json.Marshal(ar)
	if err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hub.broadcast <- jsAr

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(jsAr)
	if err != nil {
		fmt.Println("error writing response: ", err)
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
		Success: success,
	}

	err = json.NewEncoder(w).Encode(res)

	return
}
