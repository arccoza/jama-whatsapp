package main

import (
	"fmt"
	"time"
	"context"
	"strings"
	"strconv"
	whatsapp "github.com/Rhymen/go-whatsapp"
	"github.com/k0kubun/pp"
)

type WhatsAppConnector struct {
	conn *whatsapp.Conn
}

func NewWhatsAppConnector(integ *Integration) *WhatsAppConnector {
	return &WhatsAppConnector{
		conn: initWhatsApp(integ, &waHandler{}),
	}
}

func (c WhatsAppConnector) Publish(pay Payload) {

}

func (c WhatsAppConnector) Subscribe(fn Handler) {

}

func (c WhatsAppConnector) Unsubscribe(fn Handler) {

}

func (c WhatsAppConnector) Query(q string) []Payload {
	return nil
}

type waHandler struct {
	conn *whatsapp.Conn
	pub func(pay Payload)
}

func (wh *waHandler) HandleError(err error) {
	fmt.Println("WhatsApp Handler error: \n", err)
}

func (wh *waHandler) HandleJsonMessage(message string) {
	fmt.Println("HandleJsonMessage\n", message)
}

func (wh *waHandler) HandleContactMessage(message whatsapp.ContactMessage) {
	fmt.Println("HandleContactMessage\n", message)
}

// func (wh *waHandler) HandleContactList(contacts []whatsapp.Contact) {
// 	fmt.Println("HandleContactList\n", contacts)
// }

// func (wh *waHandler) HandleNewContact(contact whatsapp.Contact) {
// 	fmt.Println("HandleNewContact\n", contact)
// }

// func (wh *waHandler) HandleBatteryMessage(message whatsapp.BatteryMessage) {
// 	fmt.Println("HandleBatteryMessage\n", message)
// }

func (wh *waHandler) HandleChatList(waChats []whatsapp.Chat) {
	chats := make([]Chat, 0, len(waChats))

	for _, waChat := range waChats {
		typ := GroupChat
		id := waChat.Jid
		owner := ""
		muted, _ := strconv.ParseBool(waChat.IsMuted)
		spam, _ := strconv.ParseBool(waChat.IsMarkedSpam)
		unread, _ := strconv.ParseInt(waChat.Unread, 10, 64)
		members := map[string]ChatMember {
			wh.conn.Info.Wid: {
				ID: wh.conn.Info.Wid,
				Role: "",
				Unread: int(unread),
				Muted: muted,
				Spam: spam,
			},
		}

		if !strings.Contains(id, "-") {
			typ = DirectChat
			id = wh.conn.Info.Wid + "+" + waChat.Jid
			owner = wh.conn.Info.Wid
			members[waChat.Jid] = ChatMember{
				ID: waChat.Jid,
				Role: "member",
			}
		} else {
			meta, _ := wh.conn.GetGroupMetaData(waChat.Jid)
			owner = meta.Owner

			for _, p := range meta.Participants {
				role := "member"
				if p.IsAdmin {
					role = "admin"
				}

				if _, ok := members[p.ID]; !ok {
					members[p.ID] = ChatMember{}
				}

				member := members[p.ID]
				member.ID = p.ID
				member.Role = role
				members[p.ID] = member
			}
		}

		// fmt.Println("META: \n", meta)

		chat := Chat{
			ID: id,
			Name: waChat.Name,
			Type: typ,
			Owner: owner,
			Protocol: "whatsapp",
			Members: members,
		}

		pp.Println("CHAT: \n", chat)

		chats = append(chats, chat)
	}

	// fmt.Println("HandleChatList\n", chats)
}

// func (wh *waHandler) HandleTextMessage(message whatsapp.TextMessage) {
// 	fmt.Println("HandleTextMessage\n", message)
// }

// func (wh *waHandler) HandleImageMessage(message whatsapp.ImageMessage) {
// 	fmt.Println("HandleImageMessage\n", message)
// }

// func (wh *waHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
// 	fmt.Println("HandleDocumentMessage\n", message)
// }

// func (wh *waHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
// 	fmt.Println("HandleVideoMessage\n", message)
// }

// func (wh *waHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
// 	fmt.Println("HandleAudioMessage\n", message)
// }

func initWhatsApp(integ *Integration, handler *waHandler) *whatsapp.Conn {
	wac, err := whatsapp.NewConn(2 * time.Second)
	if err != nil {
		fmt.Println("WhatsApp connect error: \n", err)
		return nil
	}

	wac.SetClientVersion(2, 2136, 10)
	wac.SetClientName("JAMA", "jama", "0,1,0")
	handler.conn = wac
	wac.AddHandler(handler)

	if integ.Whatsapp == nil {
		integ.Whatsapp = &WhatsAppIntegration{}
	}

	// Attempt to restore the session
	if integ.Whatsapp.Session.ServerToken != "" {
		if session, err := wac.RestoreWithSession(integ.Whatsapp.Session); err == nil {
			integ.Whatsapp.Session = session
			integ.ExID = session.Wid
			integ.ref.Set(context.Background(), integ)
			return wac
		}
	}

	// If you can't restore the session request auth
	qrChan := make(chan string)
	go func() {
		val := <-qrChan
		qrToTerminal(val)
		integ.Whatsapp.QRValue = val
		integ.ref.Set(context.Background(), integ)
	}()

	if session, err := wac.Login(qrChan); err == nil {
		integ.Whatsapp.Session = session
		integ.ExID = session.Wid
		integ.ref.Set(context.Background(), integ)
		return wac
	}

	return nil
}
