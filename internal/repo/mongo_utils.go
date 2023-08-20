package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var Db *mongo.Database

const (
	NUM_CONNECTION_RETRIES = 3
	CONNECTION_TIMEOUT_SEC = 20 * time.Second
	accountCollectionName  = "account"
	groupCollectionName    = "group"
	serviceCollectionName  = "service"
)

func ConnectToMongo(ctx context.Context, uri string, dbName string) {
	for i := 0; i < NUM_CONNECTION_RETRIES; i++ {
		ctxConnect, cancel := context.WithTimeout(ctx, CONNECTION_TIMEOUT_SEC)
		client, err := mongo.Connect(ctxConnect, options.Client().ApplyURI(uri))
		cancel()
		if err != nil {
			log.Println(err)
			continue
		}
		ctxPing, cancel := context.WithTimeout(ctx, 20*time.Second)
		err = client.Ping(ctxPing, nil)
		cancel()
		if err != nil {
			log.Println(err)
			continue
		}
		Db = client.Database(dbName)
		accountCollection = Db.Collection(accountCollectionName)
		groupCollection = Db.Collection(groupCollectionName)
		serviceCollection = Db.Collection(serviceCollectionName)
		return
	}
	log.Panic("Connection to mongoDb was not set")
}

func DropDb() {
	Db.Drop(context.Background())
}
