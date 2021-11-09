package main

import (
	"fmt"
	"cloud.google.com/go/firestore"
	"context"
	whatsapp "github.com/Rhymen/go-whatsapp"
)

type Bridge struct {
	A Connector
	B Connector
}

func (br Bridge) Build(a, b Connector) {

}

func (br Bridge) Close() {

}

type BridgeManager struct {
	ctx          context.Context
	integrations *firestore.CollectionGroupRef
	bridges      map[string]Bridge
}

func NewBridgeManager(ctx context.Context, db *firestore.Client) *BridgeManager {
	return &BridgeManager{
		ctx:          ctx,
		integrations: db.CollectionGroup("integrations"),
		bridges:       make(map[string]Bridge),
	}
}

func (bm *BridgeManager) Listen() {
	it := bm.integrations.Where("kind", "==", "CHAT").Snapshots(bm.ctx)
	defer it.Stop()

	for {
		snap, err := it.Next()
		if err != nil {
			fmt.Println("Snapshots.Next: \n", err)
		}
		for _, change := range snap.Changes {
			integ := &Integration{}
			change.Doc.DataTo(integ)
			integ.ID = change.Doc.Ref.ID
			integ.Org = change.Doc.Ref.Parent.Parent.ID
			integ.ref = change.Doc.Ref

			switch change.Kind {
			case firestore.DocumentAdded:
				go bm.addBridge(integ)
			case firestore.DocumentModified:
				bm.modBridge(integ)
			case firestore.DocumentRemoved:
				bm.remBridge(integ)
			}
		}
	}
}

func (bm *BridgeManager) addBridge(integ *Integration) {
	a := NewJamaConnector(context.Background(), integ, db, fs, "", "whatsapp")
	b := NewWhatsAppConnector(integ)

	b.Subscribe(a.Publish)
	// a.Subscribe(b.Publish)

	if err := b.Start(); err != nil {
		fmt.Println("Error starting Connector: \n", err)
		return
	}
	a.Start()

	bm.bridges[integ.ID] = Bridge{a, b}
	fmt.Println("Bridge Added: \n", bm.bridges[integ.ID].A, integ)
}

func (bm *BridgeManager) modBridge(integ *Integration) {
	// bridge := bm.bridges[integ.ID]
}

func (bm *BridgeManager) remBridge(integ *Integration) {
	bridge := bm.bridges[integ.ID]
	bridge.Close()
	delete(bm.bridges, integ.ID)
}

func main() {
	jc := NewBridgeManager(context.Background(), db)
	jc.Listen()
}

type Integration struct {
	ID       string `json:"-" firestore:"-"`
	Org      string `json:"-" firestore:"-"`
	Name     string `json:"name" firestore:"name"`
	Owner    string `json:"owner" firestore:"owner"`
	ExID     string `json:"exId" firestore:"exId"` // User's External ID
	InID     string `json:"inId" firestore:"inId"` // User's Internal ID
	Provider string `json:"provider" firestore:"provider"`
	Kind     string `json:"kind" firestore:"kind"`
	ref      *firestore.DocumentRef `json:"-" firestore:"-"`
	Whatsapp *WhatsAppIntegration `json:"whatsapp" firestore:"whatsapp"`
	Timestamps *Timestamps `json:"timestamps" firestore:"timestamps"`
}

type WhatsAppIntegration struct {
	QRValue  string `json:"qrValue" firestore:"qrValue"`
	Session  whatsapp.Session `json:"session" firestore:"session"`
}

type Timestamps struct {
	Created int64 `json:"created" firestore:"created"`
	Updated int64 `json:"updated" firestore:"updated"`
}
