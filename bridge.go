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
				go bm.modBridge(integ)
			case firestore.DocumentRemoved:
				go bm.remBridge(integ)
			}
		}
	}
}

func (bm *BridgeManager) addBridge(integ *Integration) {
	a := NewJamaConnector(context.Background(), integ, db, fs, "", "whatsapp")
	b := NewWhatsAppConnector(integ)

	b.Subscribe(a.Publish)

	if err := b.Start(); err != nil {
		fmt.Println("Error starting Connector: \n", err)
		return
	}

	// fmt.Println("$$$$>", b.conn)
	a.Subscribe(b.Publish)
	a.Start()

	bm.bridges[integ.ID] = Bridge{a, b}
	fmt.Println("Bridge Added: \n", bm.bridges[integ.ID].A, integ)
}

func (bm *BridgeManager) modBridge(integ *Integration) {
	if integ.Whatsapp == nil {
		bm.addBridge(integ)
	}
}

func (bm *BridgeManager) remBridge(integ *Integration) {
	bridge := bm.bridges[integ.ID]
	bridge.Close()
	delete(bm.bridges, integ.ID)
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
	Telegram *TelegramIntegration `json:"telegram" firestore:"telegram"`
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

type TelegramIntegration struct {
	Session string
}

func (ti *TelegramIntegration) LoadSession(ctx context.Context) ([]byte, error) {
	fmt.Println("LoadSession: ", ti.Session)
	return []byte(ti.Session), nil
}

func (ti *TelegramIntegration) StoreSession(ctx context.Context, data []byte) error {
	fmt.Println("StoreSession: ", string(data))
	ti.Session = string(data)
	return nil
}
