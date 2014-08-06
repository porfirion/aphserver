package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	//"sync"
	//"net/url"
	//"time"
)

const (
	HTTP_HOST string = ""
	HTTP_PORT string = "8080"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var incomingMessages chan MessageInterface = make(chan MessageInterface, 100)
var incomingConnections chan *Connection = make(chan *Connection)

var connectionsManager ConnectionsManager = NewConnectionsManager()

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

	conn := NewConnection(ws, incomingMessages)

	incomingConnections <- conn
}

func Send(uuid string, msg MessageInterface) {
	connectionsManager.GetConnectionByUUID(uuid).ws.WriteMessage(websocket.TextMessage, StringifyMessage(msg))
}

func SendAll(msg MessageInterface) {
	for _, conn := range connectionsManager.connections {
		conn.ws.WriteMessage(websocket.TextMessage, StringifyMessage(msg))
	}
}

func logic() {
	for {
		select {
		case msg := <-incomingMessages:
			fmt.Println(StringifyMessage(msg))

			if leaveMsg, ok := msg.(*LeaveMessage); ok {
				connectionsManager.RemoveConnectionByUUID(leaveMsg.Uuid)
			}

			SendAll(msg)

		case conn := <-incomingConnections:
			log.Println("Connected: " + conn.uuid)
			connectionsManager.AddConnection(conn)

			msg := &JoinMessage{Message{MessageTypeJoin, conn.uuid}}

			SendAll(msg)
			membersUuids := connectionsManager.GetConnectionsUUIDs()

			Send(conn.uuid, SynchMembersMessage{Message{MessageTypeSynchMembers, conn.uuid}, membersUuids})
		}
	}
}

func main() {
	go logic()

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
