package main

import (
	// "encoding/json"
	nanoid "github.com/matoous/go-nanoid/v2"
	"cloud.google.com/go/firestore"
	whatsapp "github.com/Rhymen/go-whatsapp"
)

type Status int

const (
	Error Status = iota
	Pending
	Acknowledged
	Received
	Accessed
)

type ContactInfo struct {
	ID string `json:"-" firestore:"-"`
	Firstname string `json:"firstname" firestore:"firstname"`
	Lastname string `json:"lastname" firestore:"lastname"`
}

type Message struct {
	ID string `json:"-" firestore:"-"`
	Timestamp uint64 `json:"timestamp" firestore:"timestamp"`
	Protocol string `json:"protocol" firestore:"protocol"` // whatsapp, wechat, google chat, FB messenger
	From string `json:"from" firestore:"from"`
	To string `json:"to" firestore:"to"`
	Text string `json:"text" firestore:"text"`
	Status Status `json:"status" firestore:"status"` // sending, sent, received, read
	Deleted bool `json:"deleted" firestore:"deleted"`
	Attachments []Attachment `json:"attachments" firestore:"attachments"`
}

func (m *Message) fromWhatsAppMessageInfo(info whatsapp.MessageInfo) {
	m.ID = info.Id
	m.Timestamp = info.Timestamp
	m.Protocol = "whatsapp"
	m.From = info.SenderJid
	m.To = info.RemoteJid
	m.Status = Status(info.Status)
}

func (m *Message) fromWhatsApp(waMsgIf interface{}) {
	switch waMsg := waMsgIf.(type) {
	case whatsapp.TextMessage:
		m.fromWhatsAppMessageInfo(waMsg.Info)
		m.Text = waMsg.Text
	default:
		// noop
	}
}

func (m *Message) toWhatsApp() interface{} {
	info := whatsapp.MessageInfo{
		Id: m.ID,
		Timestamp: m.Timestamp,
		SenderJid: m.From,
		RemoteJid: m.To,
	}

	waMsg := whatsapp.TextMessage{
		Info: info,
		Text: m.Text,
	}

	return waMsg
}

type MessageContext struct {

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
	Org      string `json:"-" firestore:"-"`
	Name     string `json:"name" firestore:"name"`
	Owner    string `json:"owner" firestore:"owner"`
	ExID     string `json:"exId" firestore:"exId"`
	Provider string `json:"provider" firestore:"provider"`
	Kind     string `json:"kind" firestore:"kind"`
	ref      *firestore.DocumentRef `json:"-" firestore:"-"`
	Whatsapp *WhatsAppIntegration `json:"whatsapp" firestore:"whatsapp"`
}

type WhatsAppIntegration struct {
	QRValue  string `json:"qrValue" firestore:"qrValue"`
	Session  whatsapp.Session `json:"session" firestore:"session"`
}
