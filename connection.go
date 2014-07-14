package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Connection struct {
	uuid   string
	ws     *websocket.Conn
	input  chan []byte
	output chan []byte
}

func (conn *Connection) reader() {
	for {
		_, message, err := conn.ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println("Message: " + string(message))
		conn.output <- []byte(message)
	}
	conn.ws.Close()
	fmt.Println("Closing websocket " + conn.ws.RemoteAddr().String())
}

func (conn *Connection) writer() {
	for {
		var message []byte = <-conn.input
		conn.ws.WriteMessage(websocket.TextMessage, message)
		fmt.Println("incoming message: " + string(message))
	}
}

func NewConnection(soc *websocket.Conn, output chan []byte) *Connection {
	var conn *Connection = &Connection{ws: soc}
	conn.output = output
	conn.input = make(chan []byte)
	go conn.reader()
	go conn.writer()
	return conn
}
