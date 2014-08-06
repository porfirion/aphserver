package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Connection struct {
	id     uint64
	uuid   string
	ws     *websocket.Conn
	output chan MessageInterface
}

func (conn *Connection) reader() {
	defer func() {
		conn.ws.Close()
		conn.output <- &LeaveMessage{Message{MessageTypeLeave, conn.uuid}}
	}()

	for {
		_, data, err := conn.ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println("Message: " + string(data))
		msg := &TextMessage{Message{MessageTypeText, conn.uuid}, string(data)}
		conn.output <- msg
	}
	fmt.Println("Closing websocket " + conn.ws.RemoteAddr().String())
}

func NewConnection(soc *websocket.Conn, output chan MessageInterface) *Connection {
	var conn *Connection = &Connection{ws: soc}
	conn.output = output

	_, uuid, err := conn.ws.ReadMessage()
	if err != nil {
		return nil
	} else {
		conn.uuid = string(uuid)
	}

	go conn.reader()
	return conn
}
