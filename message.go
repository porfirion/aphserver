package main

import (
	"encoding/json"
	"errors"
	"log"
)

const (
	MessageTypeLogin = 1
	MessageTypeJoin  = 2
	MessageTypeLeave = 3
	MessageTypeText  = 4

	MessageTypeSynchMembers = 101
)

type MessageInterface interface {
	GetType() int
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

func ParseMessage(data []byte) (MessageInterface, error) {
	var wrapper *MessageWrapper = &MessageWrapper{}
	if err := json.Unmarshal(data, wrapper); err != nil {
		log.Println("error: ", err)
		return nil, err
	}

	//log.Println("data: ", wrapper)

	var msg MessageInterface = nil

	switch wrapper.MessageType {
	case MessageTypeLogin:
		//log.Println("login message")
		msg = &LoginMessage{}
	case MessageTypeJoin:
		//log.Println("join message")
		msg = &JoinMessage{}
	case MessageTypeLeave:
		//log.Println("leave message")
		msg = &LeaveMessage{}
	case MessageTypeText:
		//log.Println("text message")
		msg = &TextMessage{}
	default:
		//log.Println("unknown message type: ", wrapper.MessageType)
		return nil, errors.New("unkwon message type")
	}

	if err := json.Unmarshal([]byte(wrapper.Data), msg); err != nil {
		return nil, err
	}
	return msg, nil
}

type MessageWrapper struct {
	MessageType int
	Data        string
}

func (msg MessageWrapper) GetData() string {
	return msg.Data
}

type LoginMessage struct {
	UUID string
}

func (msg LoginMessage) GetType() int {
	return MessageTypeLogin
}

type JoinMessage struct {
	UUID string
}

func (msg JoinMessage) GetType() int {
	return MessageTypeJoin
}

type LeaveMessage struct {
	UUID string
}

func (msg LeaveMessage) GetType() int {
	return MessageTypeLeave
}

type TextMessage struct {
	Text string
}

func (msg TextMessage) GetType() int {
	return MessageTypeText
}

type SynchMembersMessage struct {
	Members []string
}

func (msg SynchMembersMessage) GetType() int {
	return MessageTypeSynchMembers
}
