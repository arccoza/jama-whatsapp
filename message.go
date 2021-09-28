package main

import (
	// "encoding/json"
	nanoid "github.com/matoous/go-nanoid/v2"
)

type Status int

const (
	Error Status = iota
	Pending
	Acknowledged
	Received
	Accessed
)

// type ChatStatus int

// const (
// 	Error ChatStatus = iota
// 	Inviting
// 	Invited
// 	Accepted
// )

type ChatType int

const (
	Direct ChatType = iota
	Group
	Bot
)

type Payload struct {
	Message *Message
	Chat *Chat
	ContactInfo *ContactInfo
}

type PayloadKind int

const (
	MessagePayload PayloadKind = iota
	ChatPayload
	UserPayload
)

type ContactInfo struct {
	ID string `json:"-" firestore:"-"`
	Firstname string `json:"firstname" firestore:"firstname"`
	Lastname string `json:"lastname" firestore:"lastname"`
}

type Chat struct {
	ID string `json:"-" firestore:"-"`
	Name string `json:"name" firestore:"name"`
	Type string `json:"type" firestore:"type"`// Group, Direct or Bot
	Protocol string `json:"protocol" firestore:"protocol"` // whatsapp, wechat, google chat, FB messenger
	Status Status `json:"status" firestore:"status"`
	Deleted bool `json:"deleted" firestore:"deleted"`
	Members map[string]Member `json:"members" firestore:"members"`
}

type Member struct {
	ID string `json:"id" firestore:"id"`
	Role string `json:"role" firestore:"role"`
	Unread int `json:"unread" firestore:"unread"`
	Muted bool `json:"muted" firestore:"muted"`
	Spam bool `json:"spam" firestore:"spam"`
}

type Message struct {
	ID string `json:"-" firestore:"-"`
	Timestamp int64 `json:"timestamp" firestore:"timestamp"`
	Protocol string `json:"protocol" firestore:"protocol"` // whatsapp, wechat, google chat, FB messenger
	From string `json:"from" firestore:"from"`
	To string `json:"to" firestore:"to"`
	Text string `json:"text" firestore:"text"`
	Status Status `json:"status" firestore:"status"` // sending, sent, received, read
	Deleted bool `json:"deleted" firestore:"deleted"`
	Attachments []Attachment `json:"attachments" firestore:"attachments"`
}

type Attachment struct {
	ID string `json:"id" firestore:"id"`
	Type int `json:"type" firestore:"type"`
	Mime string `json:"mime" firestore:"mime"`
	URL string `json:"url" firestore:"url"`
	Location [2]int64 `json:"location" firestore:"location"`
}

type File struct {
	ID string `json:"-" firestore:"-"`
	Type int `json:"type" firestore:"type"`
	Mime string `json:"mime" firestore:"mime"`
	URL string `json:"url" firestore:"url"`
}

type Handler func(pay Payload)

type Cache struct {
	Chats map[string]Chat
	Contacts map[string]ContactInfo
}

func NewCache() *Cache {
	return &Cache{
		make(map[string]Chat),
		make(map[string]ContactInfo),
	}
}

func (c *Cache) SetChats(chats []Chat) {
	for _, chat := range chats {
		c.Chats[chat.ID] = chat
	}
}

func (c *Cache) GetChat(id string) Chat {
	return c.Chats[id]
}

func (c *Cache) SetContacts(contacts []ContactInfo) {
	for _, contact := range contacts {
		c.Contacts[contact.ID] = contact
	}
}

func (c *Cache) GetContact(id string) ContactInfo {
	return c.Contacts[id]
}

func GenerateID() (string, error) {
	return nanoid.Generate("0123456789ABCDEF", 32)
}

type Integration struct {
	ID       string `json:"-" firestore:"-"`
	Name     string `json:"name" firestore:"name"`
	Owner    string `json:"owner" firestore:"owner"`
	Provider string `json:"provider" firestore:"provider"`
	Kind     string `json:"kind" firestore:"kind"`
	QRValue  string `json:"qrValue" firestore:"qrValue"`
	Session  string `json:"session" firestore:"session"`
}

type Connector interface {
	Publish(pay Payload)
	Subscribe(fn Handler)
	Query(q string) []Payload
}
