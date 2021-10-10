package main

import (
	// "encoding/json"
	// "cloud.google.com/go/firestore"
	whatsapp "github.com/Rhymen/go-whatsapp"
	"strings"
	"strconv"
)

type ChatType int

const (
	DirectChat ChatType = iota
	GroupChat
	BotChat
)

// type ChatStatus int

// const (
// 	Error ChatStatus = iota
// 	Inviting
// 	Invited
// 	Accepted
// )

type ChatProtocol int

const (
	UnknownProtocol ChatProtocol = iota
	WhatsAppProtocol
)

type Chat struct {
	ID string `json:"-" firestore:"-"`
	Name string `json:"name" firestore:"name"`
	Type ChatType `json:"type" firestore:"type"`// Group, Direct or Bot
	Owner string `json:"owner" firestore:"owner"`
	Protocol ChatProtocol `json:"protocol" firestore:"protocol"` // whatsapp, wechat, google chat, FB messenger
	Status Status `json:"status" firestore:"status"`
	Deleted bool `json:"deleted" firestore:"deleted"`
	Members map[string]ChatMember `json:"members" firestore:"members"`
	Timestamp int64 `json:"updated" firestore:"updated"`
}

type ChatMember struct {
	ID string `json:"id" firestore:"id"`
	Role string `json:"role" firestore:"role"`
	Unread *int `json:"unread,omitempty" firestore:"unread,omitempty"`
	Muted *bool `json:"muted,omitempty" firestore:"muted,omitempty"`
	Spam *bool `json:"spam,omitempty" firestore:"spam,omitempty"`
}

func (c *Chat) fromWhatsApp(waChat whatsapp.Chat, wac *whatsapp.Conn) error {
	cid := NormalizeWhatsAppId(waChat.Jid) // Chat id
	uid := wac.Info.Wid // User id
	mid := cid // Member id
	timestamp, _ := strconv.Atoi(waChat.LastMessageTime)

	c.Name = waChat.Name
	c.Protocol = WhatsAppProtocol
	c.Timestamp = int64(timestamp)

	muted, _ := strconv.ParseBool(waChat.IsMuted)
	spam, _ := strconv.ParseBool(waChat.IsMarkedSpam)
	unread, _ := strconv.Atoi(waChat.Unread)

	c.Members = map[string]ChatMember {
		uid: {
			ID: uid,
			Role: "owner",
			Unread: &unread,
			Muted: &muted,
			Spam: &spam,
		},
	}

	// If it's a direct chat
	if !strings.Contains(cid, "@g.us") {
		c.Type = DirectChat
		c.ID = genChatId(int(WhatsAppProtocol), int(DirectChat), []string{uid, cid})
		c.Owner = uid

		c.Members[mid] = ChatMember{
			ID: mid,
			Role: "member",
		}
	} else { // If it's a group chat
		c.Type = GroupChat
		c.ID = genChatId(int(WhatsAppProtocol), int(GroupChat), strings.Split(cid, "-"))

		if meta, err := wac.GetGroupMetaData(waChat.Jid); err != nil {
			return err
		} else {
			c.Owner = NormalizeWhatsAppId(meta.Owner)

			for _, p := range meta.Participants {
				mid = NormalizeWhatsAppId(p.ID)

				if _, ok := c.Members[mid]; !ok {
					c.Members[mid] = ChatMember{}
				}

				member := c.Members[mid]
				member.ID = mid
				member.Role = "member"

				if p.IsAdmin {
					member.Role = "admin"
				}

				c.Members[mid] = member
			}
		}
	}

	return nil
}

func genChatId(prot, typ int, parts []string) string {
	nums := make([]int, 0, 4)
	nums = append(nums, prot, typ)

	for _, part := range parts {
		if num, err := strconv.Atoi(StripWhatsAppAt(part)); err != nil {
			panic(err)
		} else {
			nums = append(nums, num)
		}
	}

	return HashIDs(nums)
}
