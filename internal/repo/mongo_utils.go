package repo

import (
	"AuthServer/internal/utils"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"time"
)

var (
	Db      *mongo.Database
	session mongo.Session
	client  *mongo.Client
)

const (
	NUM_CONNECTION_RETRIES = 3
	CONNECTION_TIMEOUT_SEC = 20 * time.Second
	accountCollectionName  = "account"
	groupCollectionName    = "group"
	serviceCollectionName  = "service"
	refreshCollectionName  = "refresh"
	accessCollectionName   = "accessC"
	bannedCollectionName   = "banned"
)

func ConnectToMongo(ctx context.Context, uri string, dbName string) {
	var err error
	for i := 0; i < NUM_CONNECTION_RETRIES; i++ {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
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
		session, err = client.StartSession()
		accountCollection = Db.Collection(accountCollectionName)
		groupCollection = Db.Collection(groupCollectionName)
		serviceCollection = Db.Collection(serviceCollectionName)
		refreshCollection = Db.Collection(refreshCollectionName)
		accessCollection = Db.Collection(accessCollectionName)
		bannedCollection = Db.Collection(bannedCollectionName)
		defaultRefreshTtl, err = strconv.ParseInt(utils.PMan.Get("default_refresh_ttl_ms").(string), 10, 64)
		if err != nil {
			panic(err)
		}
		defaultAccessTtl, err = strconv.ParseInt(utils.PMan.Get("default_access_ttl_ms").(string), 10, 64)
		if err != nil {
			panic(err)
		}
		defaultBanTtl, err = strconv.ParseInt(utils.PMan.Get("default_ban_ttl_ms").(string), 10, 64)
		if err != nil {
			panic(err)
		}
		SignSecret = []byte(utils.PMan.Get("HMAC_SECRET_KEY").(string))
		return
	}
	log.Panic("Connection to mongoDb was not set")
}

func DropDb() {
	accessCollection.Drop(context.TODO())
	refreshCollection.Drop(context.TODO())
	accountCollection.Drop(context.TODO())
	serviceCollection.Drop(context.TODO())
	bannedCollection.Drop(context.TODO())
	groupCollection.Drop(context.TODO())
}
