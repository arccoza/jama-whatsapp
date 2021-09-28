package main

import (
	whatsapp "github.com/Rhymen/go-whatsapp"
)

type WhatsAppConnector struct {
	conn *whatsapp.Conn
}

func NewWhatsAppConnector() {

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
