package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Connection struct {
	uuid   string
	ws     *websocket.Conn
	input  chan MessageInterface
	output chan MessageInterface
}

func (conn *Connection) reader() {
	for {
		_, data, err := conn.ws.ReadMessage()
		if err != nil {
			conn.output <- &LeaveMessage{Message{MessageTypeLeave, conn.uuid}}
			break
		}
		fmt.Println("Message: " + string(data))
		msg := &TextMessage{Message{MessageTypeText, conn.uuid}, string(data)}
		conn.output <- msg
	}
	conn.ws.Close()
	fmt.Println("Closing websocket " + conn.ws.RemoteAddr().String())
}

func (conn *Connection) writer() {
	for {
		msg := <-conn.input
		conn.ws.WriteMessage(websocket.TextMessage, StringifyMessage(msg))
		fmt.Println("incoming message: " + string(StringifyMessage(msg)))
	}
}

func NewConnection(soc *websocket.Conn, output chan MessageInterface) *Connection {
	var conn *Connection = &Connection{ws: soc}
	conn.output = output
	conn.input = make(chan MessageInterface)

	_, uuid, err := conn.ws.ReadMessage()
	if err != nil {
		return nil
	} else {
		conn.uuid = string(uuid)
	}

	go conn.reader()
	go conn.writer()
	return conn
}
