//TODO: simplify errors

package repo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var groupCollection *mongo.Collection

type GroupInfo struct {
	Name string `bson:"Name"`
}

func isGroupExists(group string) bool {
	res := groupCollection.FindOne(context.Background(), bson.D{{"Name", group}})
	if res.Err() == mongo.ErrNoDocuments {
		return false
	} else if res.Err() != nil {
		panic(res.Err())
	}
	return true
}

func CreateGroup(group string) (string, error) {
	if isGroupExists(group) {
		return "", errors.New("Group " + group + " already exists")
	}
	res, err := groupCollection.InsertOne(context.Background(), GroupInfo{group})
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func DeleteGroup(group string) error {
	_, err := groupCollection.DeleteOne(context.Background(), GroupInfo{group})
	return err
}
