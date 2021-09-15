package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// "google.golang.org/api/iterator"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var ctx = context.Background()
var cred = option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIAL_JSON")))
var conf = &firebase.Config{
	ProjectID: os.Getenv("FIREBASE_PROJECT_ID"),
}

func initFirebase(ctx context.Context) (*firebase.App, *firestore.Client, *storage.Client) {
	app, err := firebase.NewApp(ctx, conf, cred)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	store, err := app.Storage(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(app, db)
	return app, db, store
}
