package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"io"
	"log"
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
		conn.output <- &LeaveMessage{conn.uuid}
	}()

	for {
		if _, data, err := conn.ws.ReadMessage(); err == io.EOF {
			// client disconnected
			break
		} else if err != nil {
			log.Println("error while reading: ", err)
			break
		} else {
			if msg, err := ParseMessage(data); err != nil {
				log.Println("error while parsing: ", err)
				break
			} else {
				conn.output <- msg
			}
		}
	}
	//log.Println("Closing websocket " + conn.ws.RemoteAddr().String())
}

func NewConnection(soc *websocket.Conn, output chan MessageInterface) (*Connection, error) {
	var conn *Connection = &Connection{ws: soc}
	conn.output = output

	_, data, err := conn.ws.ReadMessage()
	msg, err := ParseMessage(data)
	if err != nil {
		return nil, err
	}
	//log.Println("message: ", msg)

	if loginMsg, ok := msg.(*LoginMessage); ok {
		conn.uuid = loginMsg.UUID
		go conn.reader()
		return conn, nil
	} else {
		return nil, errors.New("Not login message received while creating connection")
	}
}
