package main

import (
	"fmt"
	"time"
	"context"
	// "strings"
	// "strconv"
	whatsapp "github.com/Rhymen/go-whatsapp"
	// "github.com/Rhymen/go-whatsapp/binary/proto"
	// "github.com/k0kubun/pp"
	"sync"
	"encoding/json"
)

type WhatsAppConnector struct {
	integ *Integration
	conn *whatsapp.Conn
	subscribers map[*Handler]Handler
}

func NewWhatsAppConnector(integ *Integration) *WhatsAppConnector {
	c := &WhatsAppConnector{
		integ: integ,
		subscribers: map[*Handler]Handler{},
	}

	// c.conn = initWhatsApp(integ, &waHandler{
	// 	notify: func(pay Payload) {
	// 		c.notify(pay)
	// 	},
	// })

	return c
}

func (c WhatsAppConnector) Start() error {
	conn, err := initWhatsApp(c.integ, &waHandler{
		notify: func(pay Payload) {
			c.notify(pay)
		},
	})

	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c WhatsAppConnector) Publish(pay Payload) {
	for _, msg := range pay.Messages {
		c.conn.Send(msg.toWhatsApp())
	}
}

func (c WhatsAppConnector) Subscribe(fn Handler) {
	c.subscribers[&fn] = fn
}

func (c WhatsAppConnector) Unsubscribe(fn Handler) {

}

func (c *WhatsAppConnector) notify(pay Payload) {
	for _, fn := range c.subscribers {
		fn(pay)
	}
}

func (c WhatsAppConnector) Query(q string) []Payload {
	return nil
}

type waHandler struct {
	conn *whatsapp.Conn
	integ *Integration
	notify func(pay Payload)
}

func (wh *waHandler) ShouldCallSynchronously() bool {
	return false
}

func (wh *waHandler) HandleError(err error) {
	fmt.Println("WhatsApp Handler error: \n", err)
}

// func (wh *waHandler) HandleRawMessage(message *proto.WebMessageInfo) {

// }

func (wh *waHandler) HandleJsonMessage(message string) {
	fmt.Println("HandleJsonMessage\n", message)
}

func (wh *waHandler) HandleContactMessage(message whatsapp.ContactMessage) {
	fmt.Println("HandleContactMessage\n", message)
}

func (wh *waHandler) HandleContactList(waContacts []whatsapp.Contact) {
	// fmt.Printf("\nHandleContactList\n")
	// fmt.Printf("%+v\n", waContacts)
	contacts := make([]Contact, 0, len(waContacts))
	var wg sync.WaitGroup

	for _, waContact := range waContacts {
		// There's a contact for broadcast? Just skip it.
		if waContact.Jid == "status@broadcast" {
			continue
		}

		contact := &Contact{}
		contact.fromWhatsApp(waContact)
		contacts = append(contacts, *contact)
		ch, _ := wh.conn.GetProfilePicThumb(waContact.Jid)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for str := range ch {
				var dat map[string]interface{}
				json.Unmarshal([]byte(str), &dat)
				if url, _ := dat["eurl"]; url != nil {
					contact.WhatsApp.Avatar = url.(string)
				}
				// fmt.Println(contact)
			}
		}()
	}

	wg.Wait()
	wh.notify(Payload{Contacts: contacts})
}

// func (wh *waHandler) HandleNewContact(waContact whatsapp.Contact) {
// 	fmt.Printf("\nHandleNewContact\n")
// 	fmt.Printf("%+v\n", waContact)
// }

// func (wh *waHandler) HandleBatteryMessage(message whatsapp.BatteryMessage) {
// 	fmt.Println("HandleBatteryMessage\n", message)
// }

func (wh *waHandler) HandleChatList(waChats []whatsapp.Chat) {
	chats := make([]Chat, 0, len(waChats))

	for _, waChat := range waChats {
		chat := &Chat{UID: wh.integ.InID}
		chat.fromWhatsApp(waChat, wh.conn)
		chats = append(chats, *chat)
	}

	// fmt.Printf("\nHandleChatList\n")
	// fmt.Printf("%+v\n", waChats)

	wh.notify(Payload{Chats: chats})
}

func (wh *waHandler) HandleTextMessage(waMsg whatsapp.TextMessage) {
	msg := &Message{}
	msg.fromWhatsApp(waMsg, wh.conn)

	pay := Payload{Messages: []Message{*msg}}
	// fmt.Printf("\nHandleTextMessage\n")
	// fmt.Printf("%+v\n", waMsg)

	wh.notify(pay)
}

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

func initWhatsApp(integ *Integration, handler *waHandler) (*whatsapp.Conn, error) {
	wac, err := whatsapp.NewConn(30 * time.Second)
	if err != nil {
		fmt.Println("WhatsApp connect error: \n", err)
		return nil, err
	}

	wac.SetClientVersion(2, 2136, 10)
	wac.SetClientName("JAMA", "jama", "0,1,0")
	handler.conn = wac
	handler.integ = integ
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
			return wac, nil
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
		return wac, nil
	} else {
		integ.Whatsapp.QRValue = ""
		integ.ref.Set(context.Background(), integ)
		return nil, err
	}

	return nil, nil
}
