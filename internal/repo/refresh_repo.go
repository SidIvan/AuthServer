package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	refreshCollection *mongo.Collection
	accessCollection  *mongo.Collection
	bannedCollection  *mongo.Collection
	defaultRefreshTtl int64
	defaultBanTtl     int64
	signSecret        []byte
)

func ErrInvalidPayload(payload Payload) error {
	return errors.New(fmt.Sprintf("invalid payload\nPayload:\n\tLogin: %s\n\tService: %s", payload.Login, payload.Service))
}

type Token struct {
	Value     string `bson:"Value"`
	CreatedAt string `bson:"CreatedAt"`
	Ttl       int64  `bson:"Ttl"`
}

type Payload struct {
	Login   string
	Service string
	jwt.MapClaims
}

func CreateRefresh(payload Payload) (string, error) {
	return CreateRefreshWithCustomTtl(payload, defaultRefreshTtl)
}

// TODO: test
func CreateRefreshWithCustomTtl(payload Payload, ttl int64) (string, error) {
	if !isAccountExist(payload.Login) || !isServiceExist(payload.Service) {
		return "", ErrInvalidPayload(payload)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenValue, err := token.SignedString(signSecret)
	if err != nil {
		return "", err
	}
	_, err = refreshCollection.UpdateOne(context.Background(), bson.D{{"Value", tokenValue}},
		bson.D{{"$set", Token{
			Value:     tokenValue,
			CreatedAt: time.Now().Format(time.RFC3339),
			Ttl:       ttl,
		}}}, options.Update().SetUpsert(true))
	if err != nil {
		return "", err
	}
	return tokenValue, nil
}

func BanRefresh(value string) error {
	var token Token
	err := refreshCollection.FindOne(context.Background(), bson.D{{"Value", value}}).Decode(&token)
	if err != nil {
		return err
	}
	_, err = refreshCollection.DeleteOne(context.Background(), bson.D{{"Value", value}})
	if err != nil {
		return err
	}
	token.Ttl = defaultBanTtl
	token.CreatedAt = time.Now().Format(time.RFC3339)
	_, err = bannedCollection.InsertOne(context.Background(), token)
	return err
}
