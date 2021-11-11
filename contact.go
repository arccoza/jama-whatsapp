package main

import (
	// "encoding/json"
	// "cloud.google.com/go/firestore"
	whatsapp "github.com/Rhymen/go-whatsapp"
)

type Contact struct {
	ID string `json:"-" firestore:"-"`
	Name string `json:"name" firestore:"name"`
	Phone string `json:"name" firestore:"name"`
	WhatsApp *WhatsAppContact `json:"whatsapp" firestore:"whatsapp"`
	mask []string
}

type WhatsAppContact struct {
	ID string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
	Phone string `json:"name" firestore:"name"`
	Avatar string `json:"name" firestore:"name"`
}

func (c *Contact) fromWhatsAppContact(waContact whatsapp.Contact) {
	name := waContact.Name
	switch {
	case waContact.Short != "":
		name = waContact.Short
	case waContact.Notify != "":
		name = waContact.Notify
	}

	c.ID = genContactId(WhatsAppProtocol, waContact.Jid)
	c.Name = name
	c.Phone = StripWhatsAppAt(waContact.Jid)
	c.WhatsApp = &WhatsAppContact{
		ID: NormalizeWhatsAppId(waContact.Jid),
		Name: name,
		Phone: StripWhatsAppAt(waContact.Jid),
	}
}

func genContactId(prot int, id string) string {
	nums := make([]int, 0, 2)
	nums = append(nums, prot)

	if num, err := strconv.Atoi(StripWhatsAppAt(id)); err != nil {
		panic(err)
	} else {
		nums = append(nums, num)
	}

	return HashIDs(nums)
}
