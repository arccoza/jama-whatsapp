package main

import (
	// "encoding/json"
	// "cloud.google.com/go/firestore"
	// whatsapp "github.com/Rhymen/go-whatsapp"
)

type ChatType int

const (
	DirectChat ChatType = iota
	GroupChat
	BotChat
)

type Chat struct {
	ID string `json:"-" firestore:"-"`
	Name string `json:"name" firestore:"name"`
	Type ChatType `json:"type" firestore:"type"`// Group, Direct or Bot
	Owner string `json:"owner" firestore:"owner"`
	Protocol string `json:"protocol" firestore:"protocol"` // whatsapp, wechat, google chat, FB messenger
	Status Status `json:"status" firestore:"status"`
	Deleted bool `json:"deleted" firestore:"deleted"`
	Members map[string]ChatMember `json:"members" firestore:"members"`
}

type ChatMember struct {
	ID string `json:"id" firestore:"id"`
	Role string `json:"role" firestore:"role"`
	Unread *int `json:"unread,omitempty" firestore:"unread,omitempty"`
	Muted *bool `json:"muted,omitempty" firestore:"muted,omitempty"`
	Spam *bool `json:"spam,omitempty" firestore:"spam,omitempty"`
}
