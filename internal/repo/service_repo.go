//TODO: simplify errors

package repo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var serviceCollection *mongo.Collection

// TODO: change slices to sets
type Service struct {
	Name            string   `bson:"Name"`
	BaseUri         string   `bson:"BaseUri"`
	AllowedAccounts []string `bson:"AllowedAccounts"`
	AllowedGroups   []string `bson:"AllowedGroups"`
	Ruchkas         []Ruchka `bson:"Ruchkas"`
}

func isServiceExist(name string) bool {
	res := serviceCollection.FindOne(context.Background(), bson.D{{"Name", name}})
	if res.Err() == mongo.ErrNoDocuments {
		return false
	} else if res.Err() != nil {
		panic(res.Err())
	}
	return true
}

// TODO: test
func (s *Service) isRuchkaExist(rName string) bool {
	for _, ruchka := range s.Ruchkas {
		if ruchka.Name == rName {
			return true
		}
	}
	return false
}

// TODO: test
func (s *Service) hasPermission(login string) (bool, error) {
	account, err := FindAccount(login)
	if err != nil {
		return false, err
	}
	return s.hasPermissionAcc(&account)
}

// TODO: test
func (s *Service) hasPermissionAcc(account *AccountInfo) (bool, error) {
	for _, allowedAcc := range s.AllowedAccounts {
		if allowedAcc == account.Login {
			return true, nil
		}
	}
	for _, group := range account.Groups {
		for _, allowedGroup := range s.AllowedGroups {
			if group == allowedGroup {
				return true, nil
			}
		}
	}
	return false, nil
}

// TODO: test
func (s *Service) hasGroupPermission(group string) (bool, error) {
	for _, allowedGroup := range s.AllowedGroups {
		if group == allowedGroup {
			return true, nil
		}
	}
	return false, nil
}

// TODO: test
func (s *Service) hasPermissionToRuchka(rName string, login string) (bool, error) {
	account, err := FindAccount(login)
	if err != nil {
		return false, err
	}
	return s.hasPermissionToRuchkaAcc(rName, &account)
}

// TODO: test
func (s *Service) hasPermissionToRuchkaAcc(rName string, account *AccountInfo) (bool, error) {
	for _, ruchka := range s.Ruchkas {
		if ruchka.Name == rName {
			return ruchka.hasPermission(account), nil
		}
	}
	return false, errors.New(fmt.Sprintf("ruchka \"%s\" does not exist", rName))
}

func CreateService(name string, baseUri string) (string, error) {
	if isServiceExist(name) {
		return "", errors.New("service \"" + name + "\" already exists")
	}
	res, err := serviceCollection.InsertOne(context.Background(), Service{
		Name:    name,
		BaseUri: baseUri,
	})
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func FindService(name string) (*Service, error) {
	if !isServiceExist(name) {
		return nil, errors.New("Service \"" + name + "\" does not exist")
	}
	res := serviceCollection.FindOne(context.Background(), bson.D{
		{"Name", name},
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var serviceInfo Service
	if res.Decode(&serviceInfo) != nil {
		return nil, errors.New("decoding failed")
	}
	return &serviceInfo, nil
}

// TODO: test
func PutService(sName string, serviceInfo *Service) error {
	if !isServiceExist(sName) {
		return errors.New("service \"" + sName + "\"does not exists")
	}
	_, err := serviceCollection.UpdateOne(context.Background(), bson.D{{"Name", sName}}, bson.D{{"$set", serviceInfo}})
	return err
}

// TODO: test
func DeleteService(sName string) error {
	service, err := FindService(sName)
	if err != nil {
		return err
	}
	err = session.StartTransaction()
	if err != nil {
		return err
	}
	err = mongo.WithSession(context.TODO(), session, func(sc mongo.SessionContext) error {
		for _, ruchka := range service.Ruchkas {
			err = DeleteRuchka(sName, ruchka.Name)
			if err != nil {
				return err
			}
		}
		_, err := serviceCollection.DeleteOne(context.Background(), bson.D{{"Name", sName}})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		errAbort := session.AbortTransaction(context.TODO())
		if errAbort != nil {
			panic(errAbort)
		}
		return err
	}
	err = session.CommitTransaction(context.TODO())
	return err
}

// TODO: test
func AddAllowedAccount(sName string, login string) error {
	if !isAccountExist(login) {
		return errors.New(fmt.Sprintf("account \"%s\" does not exist", login))
	}
	if !isServiceExist(sName) {
		return errors.New(fmt.Sprintf("service \"%s\" does not exist", sName))
	}
	service, err := FindService(sName)
	if res, _ := service.hasPermission(login); res {
		return nil
	}
	if err != nil {
		return err
	}
	service.AllowedAccounts = append(service.AllowedAccounts, login)
	return PutService(sName, service)
}

// TODO: test
func DeleteAllowedAccount(sName string, login string) error {
	if !isAccountExist(login) {
		return errors.New(fmt.Sprintf("account \"%s\" does not exist", login))
	}
	if !isServiceExist(sName) {
		return errors.New(fmt.Sprintf("service \"%s\" does not exist", sName))
	}
	service, err := FindService(sName)
	if err != nil {
		return err
	}
	for i := range service.AllowedAccounts {
		if service.AllowedAccounts[i] == login {
			service.AllowedAccounts[i] = service.AllowedAccounts[len(service.AllowedAccounts)-1]
			service.AllowedAccounts = service.AllowedAccounts[:len(service.AllowedAccounts)-1]
			return nil
		}
	}
	return errors.New(fmt.Sprintf("account \"%s\" already has no permitions to service \"%s\"", login, sName))
}

// TODO: test
func AddAllowedGroup(sName string, group string) error {
	if !isGroupExists(group) {
		return errors.New(fmt.Sprintf("group \"%s\" does not exist", group))
	}
	if !isServiceExist(sName) {
		return errors.New(fmt.Sprintf("service \"%s\" does not exist", sName))
	}
	service, err := FindService(sName)
	if err != nil {
		return err
	}
	if res, _ := service.hasGroupPermission(group); res {
		return nil
	}
	service.AllowedGroups = append(service.AllowedGroups, group)
	return PutService(sName, service)
}

// TODO: test
func DeleteAllowedGroup(sName string, group string) error {
	if !isGroupExists(group) {
		return errors.New(fmt.Sprintf("group \"%s\" does not exist", group))
	}
	if !isServiceExist(sName) {
		return errors.New(fmt.Sprintf("service \"%s\" does not exist", sName))
	}
	service, err := FindService(sName)
	if err != nil {
		return err
	}
	for i := range service.AllowedGroups {
		if service.AllowedGroups[i] == group {
			service.AllowedGroups[i] = service.AllowedGroups[len(service.AllowedGroups)-1]
			service.AllowedGroups = service.AllowedGroups[:len(service.AllowedGroups)-1]
			return nil
		}
	}
	return errors.New(fmt.Sprintf("group \"%s\" already has no permitions to service \"%s\"", group, sName))
}
