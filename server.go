package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	//"net/url"
	"time"
)

const (
	HTTP_HOST string = ""
	HTTP_PORT string = "8080"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connections map[string]*Connection = make(map[string]*Connection)

func indexHandler(rw http.ResponseWriter, request *http.Request) {
	var indexTempl = template.Must(template.ParseFiles("templates/index.html"))
	data := struct{}{}
	indexTempl.Execute(rw, data)
}

func wsHandler(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(rw, r, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(rw, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	conn := NewConnection(ws)

	timestamp := time.Now().Format("01-01-2011 00:00:00")

	fmt.Println(timestamp)
	connections[timestamp] = conn
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ws", wsHandler)

	http.HandleFunc("/assets/", func(rw http.ResponseWriter, request *http.Request) {
		http.ServeFile(rw, request, request.URL.Path[1:])
	})

	log.Println("ADDR: " + HTTP_HOST + ":" + HTTP_PORT)

	if err := http.ListenAndServe(HTTP_HOST+":"+HTTP_PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	log.Println("running")

	for {

	}
}
