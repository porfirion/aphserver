package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
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

	if conn, err := NewConnection(ws, incomingMessages); err != nil {
		log.Println("Error on connection: ", err)
		return
	} else {
		incomingConnections <- conn
	}
}

func Send(uuid string, msg MessageInterface) {
	if data, err := StringifyMessage(msg); err == nil {
		connectionsManager.GetConnectionByUUID(uuid).ws.WriteMessage(websocket.TextMessage, data)
	}
}

func SendAll(msg MessageInterface) {
	if data, err := StringifyMessage(msg); err == nil {
		for _, conn := range connectionsManager.connections {
			conn.ws.WriteMessage(websocket.TextMessage, data)
		}
	}
}

// func SendPrecise(targets []string, exclude []string, msg MessageInterface) error {
// 	targets = sort.Strings(targets)
// 	exclude = sort.String(exclude)

// 	data, err := StringifyMessage(msg)
// 	if err != nil {
// 		return err
// 	}

// 	if len(targets) == 1 {
// 		// personal message
// 		conn := connectionsManager.GetConnectionByUUID(targets[0])
// 		conn.ws.WriteMessage(websocket.TextMessage, data)
// 	} else {

// 	}
// }

func logic() {
	for {
		select {
		case msg := <-incomingMessages:

			switch thisMsg := msg.(type) {
			case *LeaveMessage:
				connectionsManager.RemoveConnectionByUUID(thisMsg.UUID)
				log.Println(thisMsg.UUID + " leaved")
				SendAll(thisMsg)
			default:
				SendAll(thisMsg)
			}

		case conn := <-incomingConnections:
			id := connectionsManager.AddConnection(conn)
			Send(conn.uuid, &WelcomeMessage{id})

			log.Println(conn.uuid + " connected")

			SendAll(&JoinMessage{conn.uuid, id})

			Send(conn.uuid, &SyncMembersMessage{connectionsManager.GetConnectionsUUIDs()})
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
