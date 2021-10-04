package main

import (
	// "encoding/json"
	// "cloud.google.com/go/firestore"
	// whatsapp "github.com/Rhymen/go-whatsapp"
)

type ContactInfo struct {
	ID string `json:"-" firestore:"-"`
	Firstname string `json:"firstname" firestore:"firstname"`
	Lastname string `json:"lastname" firestore:"lastname"`
}
