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
	defaultAccessTtl  int64
	defaultBanTtl     int64
	SignSecret        []byte
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
	Ruchka  string
	jwt.MapClaims
}

type TokenType int

const (
	Access  TokenType = iota
	Refresh TokenType = iota
)

// TODO: test
func CreateAccess(payload Payload) (string, error) {
	return CreateToken(payload, Access)
}

// TODO: test
func CreateRefresh(payload Payload) (string, error) {
	return CreateToken(payload, Refresh)
}

func CreateToken(payload Payload, tokenType TokenType) (string, error) {
	if tokenType == Refresh {
		return CreateTokenWithCustomTtl(payload, defaultRefreshTtl, tokenType)
	}
	return CreateTokenWithCustomTtl(payload, defaultAccessTtl, tokenType)
}

// TODO: test
func CreateAccessWithCustomTtl(payload Payload) (string, error) {
	return CreateTokenWithCustomTtl(payload, defaultAccessTtl, Access)
}

// TODO: test
func CreateRefreshWithCustomTtl(payload Payload) (string, error) {
	return CreateTokenWithCustomTtl(payload, defaultRefreshTtl, Refresh)
}

// TODO: test
func CreateTokenWithCustomTtl(payload Payload, ttl int64, tokenType TokenType) (string, error) {
	collection := refreshCollection
	if tokenType == Access {
		collection = accessCollection
	}
	fmt.Println(collection)
	if !isAccountExist(payload.Login) || !isServiceExist(payload.Service) {
		return "", ErrInvalidPayload(payload)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenValue, err := token.SignedString(SignSecret)
	if err != nil {
		return "", err
	}
	_, err = collection.UpdateOne(context.Background(), bson.D{{"Value", tokenValue}},
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

// TODO: test
func BanAccess(value string) error {
	return BanToken(value, Access)
}

// TODO: test
func BanRefresh(value string) error {
	return BanToken(value, Refresh)
}

// TODO: test
func BanToken(value string, tokenType TokenType) error {
	collection := refreshCollection
	if tokenType == Access {
		collection = accessCollection
	}
	var token Token
	err := collection.FindOne(context.Background(), bson.D{{"Value", value}}).Decode(&token)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(context.Background(), bson.D{{"Value", value}})
	if err != nil {
		return err
	}
	token.Ttl = defaultBanTtl
	token.CreatedAt = time.Now().Format(time.RFC3339)
	_, err = bannedCollection.InsertOne(context.Background(), token)
	return err
}

// TODO: test
func DeleteRefresh(value string) error {
	return DeleteToken(value, Refresh)
}

// TODO: test
func DeleteAccess(value string) error {
	return DeleteToken(value, Access)
}

// TODO: test
func DeleteToken(value string, tokenType TokenType) error {
	collection := refreshCollection
	if tokenType == Access {
		collection = accessCollection
	}
	var err error
	for i := 0; i < 3; i++ {
		_, err = collection.DeleteOne(context.Background(), bson.D{{"Value", value}})
		if err == nil {
			return nil
		}
	}
	return err
}

// TODO: test
func ParseToken(tokenValue string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenValue, &Payload{}, func(t *jwt.Token) (interface{}, error) {
		return SignSecret, nil
	})
}

// TODO: test
func IsAccess(tokenValue string) bool {
	res := accessCollection.FindOne(context.Background(), bson.D{{"Value", tokenValue}})
	if res.Err() == mongo.ErrNoDocuments {
		return false
	}
	return true
}

// TODO: test
func IsRefresh(tokenValue string) bool {
	res := refreshCollection.FindOne(context.Background(), bson.D{{"Value", tokenValue}})
	if res.Err() == mongo.ErrNoDocuments {
		return false
	}
	return true
}
