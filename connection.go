package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Connection struct {
	uuid string
	ws   *websocket.Conn
}

func (conn *Connection) reader() {
	for {
		_, message, err := conn.ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println(string(message))
	}
	conn.ws.Close()
	fmt.Println("Closing websocket " + conn.ws.RemoteAddr().String())
}

func (conn *Connection) writer() {

}

func NewConnection(soc *websocket.Conn) *Connection {
	var conn *Connection = &Connection{ws: soc}
	go conn.reader()
	return conn
}
