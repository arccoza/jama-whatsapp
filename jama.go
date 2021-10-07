package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/storage"
	"fmt"
	// "google.golang.org/api/iterator"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

type JamaConnector struct {
	ctx         context.Context
	integ       *Integration
	db          *firestore.Client
	fs       *storage.Client
	uid         string
	protocol    string
	contacts    *firestore.CollectionRef
	chats       *firestore.CollectionRef
	messages    *firestore.CollectionRef
	subscribers map[*Handler]Handler
	cache       *Cache
}

func NewJamaConnector(
		ctx context.Context,
		integ *Integration,
		db *firestore.Client,
		fs *storage.Client,
		uid,
		protocol string,
	) *JamaConnector {

	return &JamaConnector{
		ctx,
		integ,
		db,
		fs,
		uid,
		protocol,
		db.Collection(fmt.Sprintf("organizations/%s/contacts", integ.Org)),
		db.Collection(fmt.Sprintf("organizations/%s/chats", integ.Org)),
		db.Collection(fmt.Sprintf("organizations/%s/messages", integ.Org)),
		map[*Handler]Handler{},
		NewCache(),
	}
}

func (c *JamaConnector) Publish(pay Payload) {
	for _, chat := range pay.Chats {
		// c.chats.Doc(chat.ID).Set(c.ctx, chat, firestore.MergeAll)
		c.chats.Doc(chat.ID).Set(c.ctx, chat)
	}

	for _, msg := range pay.Messages {
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

func (c *JamaConnector) Start() {
	// cit := c.chats.Where("status", "==", Pending).Snapshots(c.ctx)
	// mit := c.messages.Where("status", "==", Pending).Snapshots(c.ctx)
	qm := c.messages.Where("status", "==", Pending).Where("origin", "==", Internal)

	go c.listen(qm, func(change firestore.DocumentChange){
		switch change.Kind {
		case firestore.DocumentAdded:
			msg := &Message{}
			change.Doc.DataTo(msg)
			fmt.Println("Added: %v\n", msg)
		case firestore.DocumentModified:
			fmt.Println("Modified: %v\n", change.Doc.Data())
		case firestore.DocumentRemoved:
			fmt.Println("Removed: %v\n", change.Doc.Data())
		}
	})
}

func (c *JamaConnector) listen(q firestore.Query, fn func(change firestore.DocumentChange)) error {
	it := q.Snapshots(c.ctx)
	defer it.Stop()

	for {
		snap, err := it.Next()
		if err != nil {
			return fmt.Errorf("Snapshots.Next: %v", err)
		}
		for _, change := range snap.Changes {
			fn(change)
		}
	}
}

func (c *JamaConnector) Query(q string) []Payload {
	return nil
}

// func main() {
// 	jc := NewJamaConnector(context.Background(), "mako_co", db, store, "", "whatsapp")
// 	jc.Publish(Payload{Message: &Message{ID: "1234", Text: "Oi Bob.", Status: Pending}})
// 	jc.Publish(Payload{Chat: &Chat{ID: "1234", Members: map[string]Member{"1": Member{ID: "1"}, "2": Member{ID: "2"}}, Status: Pending}})
// 	jc.Listen()
// }

// type ChangeIterator struct {
// 	query   firestore.QuerySnapshotIterator
// 	idx int
// }

// func NewChangeIterator(query firestore.QuerySnapshotIterator) *ChangeIterator {
// 	return &ChangeIterator{query}
// }

// func (it *ChangeIterator) Next() (*firestore.DocumentChange, error) {
// 	snap, err := it.query.Next()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if it.idx < len(snap.Changes) {
// 		chg := snap.Changes[it.idx]
// 		it.idx += 1
// 		return &chg, nil
// 	}
// 	return nil, iterator.Done
// }

// func (it *ChangeIterator) Stop() {
// 	it.query.Stop()
// }

// type DocumentIterator struct {
// 	query    firestore.QuerySnapshotIterator
// 	docs     firestore.DocumentIterator
// 	docsDone bool
// }

// func NewDocumentIterator(query firestore.QuerySnapshotIterator) *DocumentIterator {
// 	return &DocumentIterator{query: query, docsDone: true}
// }

// func (it *DocumentIterator) Next() (*firestore.DocumentSnapshot, error) {
// 	if !it.docsDone {
// 		if doc, err := it.docs.Next; err != nil {
// 			it.docsDone = true
// 		} else {
// 			return doc, nill
// 		}
// 	}
// 	if it.docsDone {
// 		if snap, err := it.query.Next(); err != nil {
// 			return nil, err
// 		} else {
// 			it.docs = snap.Documents
// 			it.docsDone = false
// 		}
// 	}
// }

// func (it *DocumentIterator) Stop() {
// 	it.docs.Stop()
// 	it.query.Stop()
// }
