package main

import (
	"fmt"
	whatsapp "github.com/Rhymen/go-whatsapp"
	"os"
	"time"
)

func main() {
	wac, err1 := whatsapp.NewConn(20 * time.Second)
	if err1 != nil {
		fmt.Println("err1")
		fmt.Println(err1)
		return
	}

	wac.AddHandler(myHandler{})

	qrChan := make(chan string)
	go func() {
		// fmt.Printf("qr code: %v\n", qrToURI(qrChan))
		qrToTerminal(qrChan)
	}()

	// fmt.Println(whatsapp.CheckCurrentServerVersion())
	// wac.SetClientName("And1", "and1", "2,2126,14")
	wac.SetClientVersion(2, 2126, 14)

	sess, err2 := wac.Login(qrChan)
	if err2 != nil {
		fmt.Println("err2")
		fmt.Println(err2)
		return
	}

	jsonSess, _ := sess.ToJSON()
	fmt.Println("Session\n", jsonSess)

	// text := whatsapp.TextMessage{
	// 	Info: whatsapp.MessageInfo{
	// 		// Id: "s1Oy9hFHDJomezqNJegP8haTwz9tNJE9",
	// 		RemoteJid: "16178690884@s.whatsapp.net",
	// 	},
	// 	Text: "Test update",
	// }

	// id, err := wac.Send(text)

	// fmt.Println("Sent\n", text, id, err)

	fmt.Println("Waiting...")
	fmt.Scanln()
}

type myHandler struct{}

func (myHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
}

func (myHandler) HandleTextMessage(message whatsapp.TextMessage) {
	fmt.Println("HandleTextMessage\n", message)
}

func (myHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	fmt.Println("HandleImageMessage\n", message)
}

func (myHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	fmt.Println("HandleDocumentMessage\n", message)
}

func (myHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	fmt.Println("HandleVideoMessage\n", message)
}

func (myHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
	fmt.Println("HandleAudioMessage\n", message)
}

func (myHandler) HandleJsonMessage(message string) {
	fmt.Println("HandleJsonMessage\n", message)
}

func (myHandler) HandleContactMessage(message whatsapp.ContactMessage) {
	fmt.Println("HandleContactMessage\n", message)
}

func (myHandler) HandleBatteryMessage(message whatsapp.BatteryMessage) {
	fmt.Println("HandleBatteryMessage\n", message)
}

func (myHandler) HandleNewContact(contact whatsapp.Contact) {
	fmt.Println("HandleNewContact\n", contact)
}
