package repo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var accountCollection *mongo.Collection

type AccountInfo struct {
	Login    string   `bson:"Login"`
	Password string   `bson:"Password"`
	Groups   []string `bson:"Groups"`
}

func isAccountExists(login string) bool {
	res := accountCollection.FindOne(context.Background(), bson.D{{"Login", login}})
	if res.Err() == mongo.ErrNoDocuments {
		return false
	} else if res.Err() != nil {
		panic(res.Err())
	}
	return true
}

func CreateAccount(login string, password string) (string, error) {
	if isAccountExists(login) {
		return "", errors.New("Login " + login + " already using")
	}
	res, err := accountCollection.InsertOne(context.Background(), AccountInfo{
		Login:    login,
		Password: password,
		Groups:   []string{},
	})
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func GiveAccountGroup(login string, group string) error {
	account, err := FindAccount(login)
	if err != nil {
		return err
	}
	for i := 0; i < len(account.Groups); i++ {
		if account.Groups[i] == group {
			return errors.New("account already has group")
		}
	}
	account.Groups = append(account.Groups, group)
	return nil
}

func DeleteAccount(login string) error {
	_, err := accountCollection.DeleteOne(context.Background(), bson.D{
		{"Login", login},
	})
	return err
}

func FindAccount(login string) (AccountInfo, error) {
	if !isAccountExists(login) {
		return AccountInfo{}, errors.New("Account with login \"" + login + "\" does not exist")
	}
	res := accountCollection.FindOne(context.Background(), bson.D{
		{"Login", login},
	})
	if res.Err() != nil {
		return AccountInfo{}, res.Err()
	}
	var accountInfo AccountInfo
	if res.Decode(&accountInfo) != nil {
		return AccountInfo{}, errors.New("decoding failed")
	}
	return accountInfo, nil
}
