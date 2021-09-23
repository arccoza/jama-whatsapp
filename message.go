package main

import (
	// "encoding/json"
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
	Contact *Contact
}

type Contact struct {
	ID string `json:"-" firestore:"-"`
	Firstname string `json:"firstname" firestore:"firstname"`
	Lastname string `json:"lastname" firestore:"lastname"`
}

type Chat struct {
	ID string `json:"-" firestore:"-"`
	UID string `json:"uid" firestore:"uid"`
	Name string `json:"name" firestore:"name"`
	Type string `json:"type" firestore:"type"`// Group, Direct or Bot
	Protocol string `json:"protocol" firestore:"protocol"` // whatsapp, wechat, google chat, FB messenger
	Status Status `json:"status" firestore:"status"`
	IsMuted bool `json:"isMuted" firestore:"isMuted"`
	IsSpam bool `json:"isSpam" firestore:"isSpam"`
}

type Message struct {
	ID string `json:"-" firestore:"-"`
	CID string `json:"cid" firestore:"cid"` // Chat ID
	Timestamp int64 `json:"timestamp" firestore:"timestamp"`
	From string `json:"from" firestore:"from"`
	To string `json:"to" firestore:"to"`
	Text string `json:"text" firestore:"text"`
	// Attachments []Attachment `json:"attachments" firestore:"attachments"`
	Status Status `json:"status" firestore:"status"` // sending, sent, received, read
}

type Attachment struct {
	ID string `json:"-" firestore:"-"`
	Type int `json:"type" firestore:"type"`
	Mime string `json:"mime" firestore:"mime"`
	URL string `json:"url" firestore:"url"`
	Location [2]int64 `json:"location" firestore:"location"`
}

type Handler func(pay Payload)

type Cache struct {
	Chats map[string]Chat
	Users map[string]User
}

func (c *Cache) SetChats(chats []Chat) {
	for _, chat := range chats {
		c.Chats[chat.ID] = chat
	}
}

func (c *Cache) GetChat(id string) Chat {
	return c.Chats[id]
}

func (c *Cache) SetUsers(users []User) {
	for _, user := range users {
		c.Users[user.ID] = user
	}
}

func (c *Cache) GetUser(id string) User {
	return c.Users[id]
}
