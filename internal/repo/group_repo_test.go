package repo

import (
	"AuthServer/internal/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func beforeGroupTest() {
	utils.PMan = utils.NewPman()
	ConnectToMongo(context.Background(), "mongodb://"+utils.PMan.Get("mongo_host").(string)+":"+utils.PMan.Get("mongo_port").(string), utils.PMan.Get("mongo_db_name").(string))
	err := groupCollection.Drop(context.Background())
	if err != nil {
		panic(err)
	}
}

func TestCreateGroup(t *testing.T) {
	beforeGroupTest()
	id, err := CreateGroup("group1")
	if err != nil || id == "" {
		t.Errorf("Group1 creation fail")
	}
	id, err = CreateGroup("group2")
	if err != nil || id == "" {
		t.Errorf("Group2 creation fail")
	}
	id, err = CreateGroup("group2")
	if err.Error() != "Group group2 already exists" || id != "" {
		t.Errorf("Wrong creation existing group message format")
	}
}

func TestIsGroupExists(t *testing.T) {
	beforeGroupTest()
	_, err := groupCollection.InsertOne(context.Background(), bson.D{
		{"Name", "group1"},
	})
	if err != nil {
		t.Errorf("Group1 creation fail")
	}
	if !isGroupExists("group1") {
		t.Errorf("Did not find group \"group1\" in DB")
	}
	if isGroupExists("group2") {
		t.Errorf("Found group \"group2\" that does not exist")
	}
}

func TestDeleteGroup(t *testing.T) {
	beforeGroupTest()
	for i := 1; i < 4; i++ {
		_, err := groupCollection.InsertOne(context.Background(), bson.D{
			{"Name", fmt.Sprintf("group%d", i)},
		})
		if err != nil {
			t.Errorf("Group%d creation fail", i)
		}
	}
	if DeleteGroup("group2") != nil {
		t.Errorf("DeleteGroup error")
	}
	for i := 1; i < 4; i++ {
		res := groupCollection.FindOne(context.Background(), bson.D{
			{"Name", fmt.Sprintf("group%d", i)},
		})
		if i == 2 && res.Err() != mongo.ErrNoDocuments {
			t.Errorf("Group2 has not deleted")
		} else if i != 2 && res.Err() == mongo.ErrNoDocuments {
			t.Errorf("Deleted more than nesessary")
		} else if i != 2 && res.Err() != nil {
			t.Errorf(res.Err().Error())
		}
	}
}
