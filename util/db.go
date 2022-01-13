package util

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
	Ctx    context.Context
}

func (db *DB) Connect() {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	db.Client = client
	db.Ctx = context.Background()
	err = db.Client.Connect(db.Ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *DB) Disconnect() {
	db.Client.Disconnect(db.Ctx)
}
