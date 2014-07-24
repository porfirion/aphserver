package main

import (
	"encoding/json"
	//"io"
	"log"
)

const (
	MessageTypeJoin  = 1
	MessageTypeLeave = 2
	MessageTypeText  = 3

	MessageTypeSynchMembers = 101
)

type MessageInterface interface {
}

type Message struct {
	Type int
	Uuid string
}

func StringifyMessage(msg MessageInterface) []byte {
	if bytes, err := json.Marshal(msg); err == nil {
		return bytes

	} else {
		log.Fatal(err)
		return nil
	}
}

type JoinMessage struct {
	Message
}

type LeaveMessage struct {
	Message
}

type TextMessage struct {
	Message
	Text string
}

type SynchMembersMessage struct {
	Message
	Members []string
}

/*
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
*/
