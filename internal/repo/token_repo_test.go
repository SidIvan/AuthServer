package repo

import (
	"AuthServer/internal/utils"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"testing"
)

func beforeRefreshTest() {
	utils.PMan = utils.NewPman("test.properties")
	ConnectToMongo(context.Background(), "mongodb://"+utils.PMan.Get("mongo_host").(string)+":"+utils.PMan.Get("mongo_port").(string), utils.PMan.Get("mongo_db_name").(string))
	err := refreshCollection.Drop(context.Background())
	if err != nil {
		panic(err)
	}
}

func TestCreateRefresh(t *testing.T) {
	beforeRefreshTest()
	CreateAccount("login", "password")
	CreateService("service", "uri")
	val, err := CreateRefresh(Payload{
		Login:   "login",
		Service: "service",
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	token, err := jwt.ParseWithClaims(val, &Payload{}, func(t *jwt.Token) (interface{}, error) {
		return SignSecret, nil
	})
	if !token.Valid {
		t.Errorf("invalid token created")
	}
	val, err = CreateRefresh(Payload{
		Login:   "login",
		Service: "service",
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	invalidLoginPayload := Payload{
		Login:   "not existing login",
		Service: "service",
	}
	_, err = CreateRefresh(invalidLoginPayload)
	if errors.Is(err, ErrInvalidPayload(invalidLoginPayload)) {
		t.Errorf("wrong error format for not existed login case")
	}
	invalidServicePayload := Payload{
		Login:   "login",
		Service: "not existing service",
	}
	_, err = CreateRefresh(invalidServicePayload)
	if errors.Is(err, ErrInvalidPayload(invalidServicePayload)) {
		t.Errorf("wrong error format for not existed service case")
	}
}
