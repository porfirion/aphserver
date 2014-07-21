package main

import (
	"encoding/json"
	"io"
	"log"
)

type Message struct {
	Type, Data string
}

func ReadMessage(buffer []byte) *Message {
	var message *Message
	if err := json.Unmarshal(buffer, message); err != nil && err != io.EOF {
		log.Fatal(err)
	}

	return message
}

func WriteMessage(msg *Message) []byte {
	var bytes []byte
	var err error
	if bytes, err = json.Marshal(msg); err != nil {
		log.Fatal(err)
	}

	return bytes
}
