package main

import (
	whatsapp "github.com/Rhymen/go-whatsapp"
)

type WhatsAppConnector struct {
	conn *whatsapp.Conn
}

func (c WhatsAppConnector) Publish(msg *Message) {

}

func (c WhatsAppConnector) Subscribe(handler *Handler) {

}

func (c WhatsAppConnector) Query() {

}
