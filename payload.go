package main

import (
	// "encoding/json"
	// "cloud.google.com/go/firestore"
	// whatsapp "github.com/Rhymen/go-whatsapp"
)

type Payload struct {
	Messages []Message
	Chats []Chat
	Contacts []Contact
}

type PayloadKind int

const (
	MessagePayload PayloadKind = iota
	ChatPayload
	UserPayload
)

type Handler func(pay Payload)

type Connector interface {
	Publish(pay Payload)
	Subscribe(fn Handler)
	Query(q string) []Payload
}
