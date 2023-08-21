package repo

import (
	"AuthServer/internal/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"testing"
)

func beforeAccountTest() {
	utils.PMan = utils.NewPman("test.properties")
	ConnectToMongo(context.Background(), "mongodb://"+utils.PMan.Get("mongo_host").(string)+":"+utils.PMan.Get("mongo_port").(string), utils.PMan.Get("mongo_db_name").(string))
	err := accountCollection.Drop(context.Background())
	if err != nil {
		panic(err)
	}
}

func TestCreateAccount(t *testing.T) {
	beforeAccountTest()
	id, err := CreateAccount("login1", "password1")
	if err != nil || id == "" {
		t.Errorf("Account1 creation fail")
	}
	id, err = CreateAccount("login2", "password2")
	if err != nil || id == "" {
		t.Errorf("Account2 creation fail")
	}
	id, err = CreateAccount("login2", "password3")
	if err != nil && err.Error() != "Login login2 already using" || id != "" {
		t.Errorf("Wrong creation existing group message format")
	}
}

func TestIsAccountExists(t *testing.T) {
	beforeAccountTest()
	_, err := accountCollection.InsertOne(context.Background(), AccountInfo{
		Login:    "login1",
		Password: "password1",
	})
	if err != nil {
		t.Errorf("Account1 creation fail")
	}
	if !isAccountExist("login1") {
		t.Errorf("Did not find account with login \"login1\" in DB")
	}
	if isAccountExist("login2") {
		t.Errorf("Found account with login \"login2\" that does not exist")
	}
}

func TestDeleteAccount(t *testing.T) {
	beforeAccountTest()
	for i := 1; i < 4; i++ {
		_, err := accountCollection.InsertOne(context.Background(), bson.D{
			{"Login", fmt.Sprintf("login%d", i)},
			{"Password", fmt.Sprintf("password%d", i)},
		})
		if err != nil {
			t.Errorf("Account%d creation fail", i)
		}
	}
	if DeleteAccount("login2") != nil {
		t.Errorf("DeleteAccount error")
	}
	for i := 1; i < 4; i++ {
		res := accountCollection.FindOne(context.Background(), bson.D{
			{"Login", fmt.Sprintf("login%d", i)},
			{"Password", fmt.Sprintf("password%d", i)},
		})
		if i == 2 && res.Err() != mongo.ErrNoDocuments {
			t.Errorf("Account2 has not deleted")
		} else if i != 2 && res.Err() == mongo.ErrNoDocuments {
			t.Errorf("Deleted more than nesessary")
		} else if i != 2 && res.Err() != nil {
			t.Errorf(res.Err().Error())
		}
	}
}

func TestFindAccount(t *testing.T) {
	beforeAccountTest()
	for i := 1; i < 4; i++ {
		_, err := accountCollection.InsertOne(context.Background(), bson.D{
			{"Login", fmt.Sprintf("login%d", i)},
			{"Password", fmt.Sprintf("password%d", i)},
			{"Groups", []string{"group1", "group2"}},
		})
		if err != nil {
			t.Errorf("Account%d creation fail", i)
		}
	}
	for i := 1; i < 4; i++ {
		res, err := FindAccount(fmt.Sprintf("login%d", i))
		if err != nil {
			t.Errorf("Error returned for existed account")
		}
		if res.Login != fmt.Sprintf("login%d", i) || res.Password != fmt.Sprintf("password%d", i) ||
			len(res.Groups) != 2 || res.Groups[0] != "group1" || res.Groups[1] != "group2" {
			t.Errorf("Wrong account found for login \"login%d\"", i)
		}
	}
}

func TestGiveAccountGroup(t *testing.T) {
	beforeAccountTest()
	_, err := accountCollection.InsertOne(context.Background(), bson.D{
		{"Login", "login"},
		{"Password", "password"},
		{"Groups", []string{"group"}},
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	err = GiveAccountGroup("login", "newGroup")
	if err != nil {
		t.Errorf(err.Error())
	} else {
		var account AccountInfo
		resultGroups := []string{"group", "newGroup"}
		if err = accountCollection.FindOne(context.Background(), bson.D{
			{"Login", "login"},
		}).Decode(&account); err != nil || reflect.DeepEqual(account.Groups, resultGroups) {
			t.Errorf(err.Error())
		}
	}
	err = GiveAccountGroup("login", "group")
	if err.Error() != "account already has group" {
		t.Errorf("error message in already existed group did not return")
	}
	err = GiveAccountGroup("nonExistenceLogin", "group")
	if err.Error() != "Account with login \"nonExistenceLogin\" does not exist" {
		t.Errorf("error message in non existence account did not return")
	}
}
