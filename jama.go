package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/storage"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JamaConnector struct {
	ctx         context.Context
	db          *firestore.Client
	store       *storage.Client
	contacts    *firestore.CollectionRef
	chats       *firestore.CollectionRef
	messages    *firestore.CollectionRef
	subscribers map[*Handler]Handler
}

func NewJamaConnector(ctx context.Context, db *firestore.Client, store *storage.Client) *JamaConnector {
	return &JamaConnector{
		ctx,
		db,
		store,
		db.Collection("contacts"),
		db.Collection("chats"),
		db.Collection("messages"),
		map[*Handler]Handler{},
	}
}

func (c *JamaConnector) Publish(pay Payload) {

	if chat := pay.Chat; chat != nil {
		c.chats.Doc(chat.ID).Set(c.ctx, chat)
	}

	if msg := pay.Message; msg != nil {
		c.messages.Doc(msg.ID).Set(c.ctx, msg)
	}
}

func (c *JamaConnector) Subscribe(fn Handler) {
	c.subscribers[&fn] = fn
}

func (c *JamaConnector) Unsubscribe(fn Handler) {
	delete(c.subscribers, &fn)
}

func (c *JamaConnector) notify(pay Payload) {
	for _, fn := range c.subscribers {
		fn(pay)
	}
}

func (c *JamaConnector) listen() error {
	it := c.messages.Where("status", "==", Pending).Snapshots(c.ctx)

	for {
		snap, err := it.Next()
		// DeadlineExceeded will be returned when ctx is cancelled.
		if status.Code(err) == codes.DeadlineExceeded {
			return nil
		}
		if err != nil {
			return fmt.Errorf("Snapshots.Next: %v", err)
		}
		if snap != nil {
			for _, change := range snap.Changes {
				switch change.Kind {
				case firestore.DocumentAdded:
					fmt.Println("Added: %v\n", change.Doc.Data())
				case firestore.DocumentModified:
					fmt.Println("Modified: %v\n", change.Doc.Data())
				case firestore.DocumentRemoved:
					fmt.Println("Removed: %v\n", change.Doc.Data())
				}
			}

			for {
				doc, err := snap.Documents.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return fmt.Errorf("Documents.Next: %v", err)
				}
				fmt.Println("Current cities in California: %v\n", doc.Ref.ID)
			}
		}

	}
}

func (c *JamaConnector) Query() {

}

func main() {
	jc := NewJamaConnector(context.Background(), db, store)
	jc.Publish(Payload{Message: &Message{ID: "1234", Text: "Oi Bob.", Status: Pending}})
	jc.listen()
}
