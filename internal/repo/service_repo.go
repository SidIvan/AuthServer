package repo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var serviceCollection *mongo.Collection

type ServiceInfo struct {
	Name            string   `bson:"Name"`
	BaseUri         string   `bson:"BaseUri"`
	AllowedAccounts []string `bson:"AllowedAccounts"`
	AllowedGroups   []string `bson:"AllowedGroups"`
	Ruchkas         []Ruchka `bson:"Ruchkas"`
}

type Ruchka struct {
	Name            string
	Uri             string
	Method          string
	AllowedAccounts []string
	AllowedGroups   []string
}

func isServiceExists(name string) bool {
	res := serviceCollection.FindOne(context.Background(), bson.D{{"Name", name}})
	if res.Err() == mongo.ErrNoDocuments {
		return false
	} else if res.Err() != nil {
		panic(res.Err())
	}
	return true
}

func CreateService(name string, baseUri string) (string, error) {
	if isServiceExists(name) {
		return "", errors.New("service \"" + name + "\" already exists")
	}
	res, err := serviceCollection.InsertOne(context.Background(), ServiceInfo{
		Name:    name,
		BaseUri: baseUri,
	})
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func FindService(name string) (*ServiceInfo, error) {
	if !isServiceExists(name) {
		return nil, errors.New("Service \"" + name + "\" does not exist")
	}
	res := serviceCollection.FindOne(context.Background(), bson.D{
		{"Name", name},
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var serviceInfo ServiceInfo
	if res.Decode(&serviceInfo) != nil {
		return nil, errors.New("decoding failed")
	}
	return &serviceInfo, nil
}

func AddRuchka(name string, ruchka Ruchka) error {
	if !isServiceExists(name) {
		return errors.New("service \"" + name + "\"does not exists")
	}
	service, err := FindService(name)
	if err != nil {
		return err
	}
	service.Ruchkas = append(service.Ruchkas, ruchka)
	_, err = serviceCollection.UpdateOne(context.Background(), bson.D{{"Name", name}}, bson.D{{"$set", service}})
	return err
}

func DeleteRuchka(name string, ruchkaName string) error {
	service, err := FindService(name)
	if err != nil {
		return err
	}
	for i := 0; i < len(service.Ruchkas); i++ {
		if service.Ruchkas[i].Name == ruchkaName {
			service.Ruchkas[i] = service.Ruchkas[len(service.Ruchkas)-1]
			service.Ruchkas = service.Ruchkas[:len(service.Ruchkas)-1]
			_, err = serviceCollection.UpdateOne(context.Background(), bson.D{{"Name", name}}, bson.D{{"$set", service}})
			return err
		}
	}
	return errors.New("ruchka \"" + ruchkaName + "\" not found")
}
