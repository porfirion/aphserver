package main

import (
	"encoding/json"
	"errors"
	"log"
)

const (
	MessageTypeLogin     = 1
	MessageTypeWelcome   = 2
	MessageTypeForbidden = 3
	MessageTypeJoin      = 10
	MessageTypeLeave     = 11

	MessageTypeSyncMembers = 101

	MessageTypeText = 1001
)

type MessageInterface interface {
	GetType() int
}

func CreateMessageByType(msgType int) (MessageInterface, error) {
	var msg MessageInterface = nil

	switch msgType {
	case MessageTypeLogin:
		msg = &LoginMessage{}
	case MessageTypeWelcome:
		msg = &WelcomeMessage{}
	case MessageTypeForbidden:
		msg = &ForbiddenMessage{}
	case MessageTypeJoin:
		msg = &JoinMessage{}
	case MessageTypeLeave:
		msg = &LeaveMessage{}
	case MessageTypeText:
		msg = &TextMessage{}
	default:
		log.Println("Unknown message type (create): ", msgType)
		return nil, errors.New("Unkwon message type (create)")
	}

	return msg, nil
}

func StringifyMessage(msg MessageInterface) ([]byte, error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	wrapper := &MessageTransportWrapper{msg.GetType(), string(bytes)}

	if bytes, err := json.Marshal(wrapper); err == nil {
		return bytes, nil
	} else {
		return nil, err
	}
}

func ParseMessage(data []byte) (MessageInterface, error) {
	var wrapper *MessageTransportWrapper = &MessageTransportWrapper{}
	if err := json.Unmarshal(data, wrapper); err != nil {
		log.Println("error: ", err)
		return nil, err
	}

	msg, err := CreateMessageByType(wrapper.MessageType)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(wrapper.Data), msg); err != nil {
		return nil, err
	}
	return msg, nil
}

type MessageTransportWrapper struct {
	MessageType int
	Data        string
}

type LoginMessage struct {
	UUID string
	Name string
}

type WelcomeMessage struct {
	Id uint64
}

type ForbiddenMessage struct {
	Cause string
}

type JoinMessage struct {
	UUID string
	Id   uint64
}

type LeaveMessage struct {
	UUID string
}

type TextMessage struct {
	Text string
}

type SyncMembersMessage struct {
	Members []string
}

func (msg LoginMessage) GetType() int {
	return MessageTypeLogin
}
func (msg WelcomeMessage) GetType() int {
	return MessageTypeWelcome
}
func (msg ForbiddenMessage) GetType() int {
	return MessageTypeForbidden
}
func (msg JoinMessage) GetType() int {
	return MessageTypeJoin
}
func (msg LeaveMessage) GetType() int {
	return MessageTypeLeave
}
func (msg TextMessage) GetType() int {
	return MessageTypeText
}
func (msg SyncMembersMessage) GetType() int {
	return MessageTypeSyncMembers
}
