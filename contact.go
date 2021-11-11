package main

import (
	// "fmt"
	"strconv"
	"strings"
	// "encoding/json"
	// "cloud.google.com/go/firestore"
	whatsapp "github.com/Rhymen/go-whatsapp"
)

type Contact struct {
	ID string `json:"-" firestore:"-"`
	Name string `json:"name" firestore:"name"`
	Phone string `json:"phone" firestore:"phone"`
	WhatsApp *WhatsAppContact `json:"whatsapp" firestore:"whatsapp"`
	mask []string
}

type WhatsAppContact struct {
	ID string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
	Phone string `json:"phone" firestore:"phone"`
	Avatar string `json:"avatar" firestore:"avatar"`
}

func (c *Contact) fromWhatsApp(waContact whatsapp.Contact) {
	name := waContact.Name
	switch {
	case waContact.Short != "":
		name = waContact.Short
	case waContact.Notify != "":
		name = waContact.Notify
	}

	c.ID = genContactId(int(WhatsAppProtocol), strings.Split(waContact.Jid, "-"))
	c.Name = name
	c.Phone = StripWhatsAppAt(waContact.Jid)
	c.WhatsApp = &WhatsAppContact{
		ID: NormalizeWhatsAppId(waContact.Jid),
		Name: name,
		Phone: StripWhatsAppAt(waContact.Jid),
	}
}

func genContactId(prot int, parts []string) string {
	nums := make([]int, 0, 3)
	nums = append(nums, prot)

	for _, part := range parts {
		if num, err := strconv.Atoi(StripWhatsAppAt(part)); err != nil {
			panic(err)
		} else {
			nums = append(nums, num)
		}
	}

	return HashIDs(nums)
}
