package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/control", serveControl)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Printf("Listening on port %s", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("serveHome called: %s", r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		log.Printf("No Found: %s", r.URL)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		log.Printf("Method not allowed: %s", r.URL)
		return
	}

	log.Printf("Serving home: %s", r.URL)
	http.ServeFile(w, r, "./html/home.html")
}

func serveControl(w http.ResponseWriter, r *http.Request) {
	log.Println("serveControl called")
	http.ServeFile(w, r, "./html/control.html")
	return
}
