//TODO: simplify errors

package repo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

// TODO: change slices to sets
type Ruchka struct {
	Name            string
	Uri             string
	Method          string
	AllowedAccounts []string
	AllowedGroups   []string
}

// TODO: test
func isRuchkaExist(sName string, rName string) bool {
	service, err := FindService(sName)
	if err != nil {
		return false
	}
	for _, ruchka := range service.Ruchkas {
		if ruchka.Name == rName {
			return true
		}
	}
	return false
}

// TODO: test
func (r *Ruchka) hasPermission(account *AccountInfo) bool {
	for _, allowedAcc := range r.AllowedAccounts {
		if allowedAcc == account.Login {
			return true
		}
	}
	return false
}

// TODO: test
func (r *Ruchka) hasGroupPermission(gName string) bool {
	for _, allowedGroup := range r.AllowedGroups {
		if allowedGroup == gName {
			return true
		}
	}
	return false
}

func AddRuchka(name string, ruchka Ruchka) error {
	if !isServiceExist(name) {
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

// TODO: test
func PutRuchka(name string, ruchkaName string, ruchka Ruchka) error {
	service, err := FindService(name)
	if err != nil {
		return err
	}
	for i := 0; i < len(service.Ruchkas); i++ {
		if service.Ruchkas[i].Name == ruchkaName {
			service.Ruchkas[i] = ruchka
			_, err = serviceCollection.UpdateOne(context.Background(), bson.D{{"Name", name}}, bson.D{{"$set", service}})
			return err
		}
	}
	return errors.New("ruchka \"" + ruchkaName + "\" not found")
}

// TODO: test
func AddRuchkaAllowedAccount(sName string, rName string, login string) error {
	if !isServiceExist(sName) {
		return errors.New(fmt.Sprintf("service \"%s\" does not exist", sName))
	}
	if !isAccountExist(login) {
		return errors.New(fmt.Sprintf("account \"%s\" does not exist", login))
	}
	service, err := FindService(sName)
	if err != nil {
		return nil
	}
	if res, err := service.hasPermissionToRuchka(rName, login); res {
		return nil
	} else if err != nil {
		return err
	}
	for i := range service.Ruchkas {
		if service.Ruchkas[i].Name == rName {
			service.Ruchkas[i].AllowedAccounts = append(service.Ruchkas[i].AllowedAccounts, login)
		}
	}
	return nil
}

// TODO: test
func DeleteRuchkaAllowedAccount(sName string, rName string, login string) error {
	if !isServiceExist(sName) {
		return errors.New(fmt.Sprintf("service \"%s\" does not exist", sName))
	}
	if !isAccountExist(login) {
		return errors.New(fmt.Sprintf("account \"%s\" does not exist", login))
	}
	service, err := FindService(sName)
	if err != nil {
		return nil
	}
	if res, err := service.hasPermissionToRuchka(rName, login); !res {
		return errors.New(fmt.Sprintf("account \"%s\" already has no permission to ruchka \"%s\" of service \"%s\"", login, rName, sName))
	} else if err != nil {
		return err
	}
	for i := range service.Ruchkas {
		if service.Ruchkas[i].Name == rName {
			service.Ruchkas[i].AllowedAccounts = append(service.Ruchkas[i].AllowedAccounts, login)
			return nil
		}
	}
	return nil
}

// TODO: test
func AddRuchkaAllowedGroup(sName string, rName string, gName string) error {
	if !isGroupExists(gName) {
		return errors.New(fmt.Sprintf("group \"%s\" does not exist", gName))
	}
	service, err := FindService(sName)
	if err != nil {
		return err
	}
	for i := range service.Ruchkas {
		if service.Ruchkas[i].Name == rName {
			if service.Ruchkas[i].hasGroupPermission(gName) {
				return nil
			}
			service.Ruchkas[i].AllowedGroups = append(service.Ruchkas[i].AllowedGroups, gName)
			err = PutService(sName, service)
			return err
		}
	}
	return errors.New(fmt.Sprintf("ruchka \"%s\" of service \"%s\" does not exist", rName, sName))
}

// TODO: test
func DeleteRuchkaAllowedGroup(sName string, rName string, gName string) error {
	if !isGroupExists(gName) {
		return errors.New(fmt.Sprintf("group \"%s\" does not exist", gName))
	}
	service, err := FindService(sName)
	if err != nil {
		return err
	}
	for i := range service.Ruchkas {
		if service.Ruchkas[i].Name == rName {
			if !service.Ruchkas[i].hasGroupPermission(gName) {
				return nil
			}
			for j, group := range service.Ruchkas[i].AllowedGroups {
				if group == gName {
					service.Ruchkas[i].AllowedGroups[j] = service.Ruchkas[i].AllowedGroups[len(service.Ruchkas[i].AllowedGroups)-1]
					service.Ruchkas[i].AllowedGroups = service.Ruchkas[i].AllowedGroups[:len(service.Ruchkas[i].AllowedGroups)-1]
					err = PutService(sName, service)
					return err
				}
			}
		}
	}
	return errors.New(fmt.Sprintf("ruchka \"%s\" of service \"%s\" does not exist", rName, sName))
}
