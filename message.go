package main

import (
	// "fmt"
	"strings"
	// "encoding/json"
	nanoid "github.com/matoous/go-nanoid/v2"
	// "cloud.google.com/go/firestore"
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

type Origin int

const (
	External Origin = iota
	Internal
)

type Message struct {
	ID string `json:"-" firestore:"-"`
	CID string `json:"cid" firestore:"cid"`
	Timestamp int64 `json:"timestamp" firestore:"timestamp"`
	Protocol string `json:"protocol" firestore:"protocol"` // whatsapp, wechat, google chat, FB messenger
	From string `json:"from" firestore:"from"`
	To string `json:"to" firestore:"to"`
	Text string `json:"text" firestore:"text"`
	Status Status `json:"status" firestore:"status"` // sending, sent, received, read
	Deleted bool `json:"deleted" firestore:"deleted"`
	Origin Origin `json:"origin" firestore:"origin"`
	Attachments []Attachment `json:"attachments" firestore:"attachments"`
}

func (m *Message) fromWhatsAppMessageInfo(info whatsapp.MessageInfo, uid string) {
	cid, from, to := "", "", ""

	if from = uid; info.FromMe {
		to = info.RemoteJid
	} else if from = info.SenderJid; info.SenderJid != "" {
		to = info.RemoteJid
	} else {
		from = info.RemoteJid
		to = uid
	}

	if IsWhatsAppGroup(info.RemoteJid) {
		cid = genChatId(1, int(GroupChat), strings.Split(info.RemoteJid, "-"))
	} else {
		cid = genChatId(1, int(DirectChat), []string{uid, info.RemoteJid})
	}

	m.ID = info.Id
	m.CID = cid
	m.Timestamp = int64(info.Timestamp)
	m.Protocol = "whatsapp"
	m.From = NormalizeWhatsAppId(from)
	m.To = NormalizeWhatsAppId(to)
	m.Status = Status(info.Status)
	// fmt.Printf("\nfromWhatsApp\n%+v\n", m)
	// fmt.Printf("\nfromWhatsApp\n %q %q %q %q %q %q %q", cid, uid, from, to, info.FromMe, info.SenderJid, info.RemoteJid)
}

func (m *Message) fromWhatsApp(waMsgIf interface{}, wac *whatsapp.Conn) {
	switch waMsg := waMsgIf.(type) {
	case whatsapp.TextMessage:
		m.fromWhatsAppMessageInfo(waMsg.Info, wac.Info.Wid)
		m.Text = waMsg.Text
	default:
		// noop
	}
}

func (m *Message) toWhatsApp() interface{} {
	info := whatsapp.MessageInfo{
		Id: m.ID,
		Timestamp: uint64(m.Timestamp),
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
