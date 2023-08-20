package repo

import (
	"AuthServer/internal/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
)

func beforeServiceTest() {
	utils.PMan = utils.NewPman()
	ConnectToMongo(context.Background(), "mongodb://"+utils.PMan.Get("mongo_host").(string)+":"+utils.PMan.Get("mongo_port").(string), utils.PMan.Get("mongo_db_name").(string))
	err := serviceCollection.Drop(context.Background())
	if err != nil {
		panic(err)
	}
}

func TestIsServiceExists(t *testing.T) {
	beforeServiceTest()
	_, err := serviceCollection.InsertOne(context.Background(), ServiceInfo{Name: "service1"})
	if err != nil {
		t.Errorf("Service1 creation fail")
	}
	if !isServiceExists("service1") {
		t.Errorf("Did not find service \"service1\" in DB")
	}
	if isServiceExists("service2") {
		t.Errorf("Found service \"service2\" that does not exist")
	}
}

func TestCreateService(t *testing.T) {
	beforeServiceTest()
	id, err := CreateService("service1", "baseUri1")
	if id == "" || err != nil {
		t.Errorf(err.Error())
	}
	id, err = CreateService("service2", "baseUri2")
	if id == "" || err != nil {
		t.Errorf(err.Error())
	}
	id, err = CreateService("service1", "baseUri3")
	if id != "" || (err != nil && err.Error() != "service \"service1\" already exists") {
		t.Errorf("wrong already exists service return format")
	}
	var service ServiceInfo
	err = serviceCollection.FindOne(context.Background(), bson.D{{"Name", "service1"}}).Decode(&service)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(service, ServiceInfo{
		Name:    "service1",
		BaseUri: "baseUri1",
	}) {
		t.Errorf("saved info does not match expected")
	}
	err = serviceCollection.FindOne(context.Background(), bson.D{{"Name", "service2"}}).Decode(&service)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(service, ServiceInfo{
		Name:    "service2",
		BaseUri: "baseUri2",
	}) {
		t.Errorf("saved info does not match expected")
	}
}

func TestFindService(t *testing.T) {
	beforeServiceTest()
	services := make(map[string]ServiceInfo)
	for i := 1; i < 4; i++ {
		serviceName := fmt.Sprintf("service%d", i)
		services[serviceName] = ServiceInfo{
			Name:            serviceName,
			BaseUri:         fmt.Sprintf("baseUri%d", i),
			AllowedAccounts: []string{fmt.Sprintf("acc1_%d", i), fmt.Sprintf("acc2_%d", i)},
			AllowedGroups:   []string{fmt.Sprintf("gr1_%d", i), fmt.Sprintf("gr2_%d", i)},
			Ruchkas: []Ruchka{
				{
					Name:            fmt.Sprintf("ruchka1_%d", i),
					Uri:             fmt.Sprintf("uri1_%d", i),
					Method:          fmt.Sprintf("method1_%d", i),
					AllowedAccounts: []string{fmt.Sprintf("acc1_%d", i), fmt.Sprintf("acc2_%d", i)},
					AllowedGroups:   []string{fmt.Sprintf("gr1_%d", i), fmt.Sprintf("gr2_%d", i)},
				},
				{
					Name:   fmt.Sprintf("ruchka2_%d", i),
					Uri:    fmt.Sprintf("uri2_%d", i),
					Method: fmt.Sprintf("method2_%d", i),
				},
			}}
		_, err := serviceCollection.InsertOne(context.Background(), services[serviceName])
		if err != nil {
			t.Errorf("Service%d creation fail", i)
		}
	}
	for serviceName, serviceInfo := range services {
		res, err := FindService(serviceName)
		if err != nil {
			t.Errorf("Error returned for existed service")
		}
		if !reflect.DeepEqual(*res, serviceInfo) {
			t.Errorf("Wrong service found for service \"" + serviceName + "\"")
		}
	}
}

func TestAddRuchka(t *testing.T) {
	beforeServiceTest()
	service := ServiceInfo{
		Name:            "serviceName",
		BaseUri:         "baseUri",
		AllowedAccounts: []string{"acc1", "acc2"},
		AllowedGroups:   []string{"gr1", "gr2"},
		Ruchkas: []Ruchka{
			{
				Name:            "ruchka1",
				Uri:             "uri1",
				Method:          "method1",
				AllowedAccounts: []string{"acc1", "acc2"},
				AllowedGroups:   []string{"gr1", "gr2"},
			},
			{
				Name:   "ruchka2",
				Uri:    "uri2",
				Method: "method2",
			},
		}}
	_, err := serviceCollection.InsertOne(context.Background(), service)
	if err != nil {
		t.Errorf(err.Error())
	}
	newRuchka := Ruchka{
		Name:            "newRuchka",
		Uri:             "newRuchkaUri",
		Method:          "newRuchkaMethod",
		AllowedAccounts: []string{"newRuchkaAllowedAccount"},
		AllowedGroups:   []string{"newRuchkaAllowedGroups"},
	}
	service.Ruchkas = append(service.Ruchkas, newRuchka)
	err = AddRuchka(service.Name, newRuchka)
	if err != nil {
		t.Errorf(err.Error())
	}
	var res ServiceInfo
	err = serviceCollection.FindOne(context.Background(), bson.D{{"Name", service.Name}}).Decode(&res)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(service, res) {
		t.Errorf("incorrect addition of rucka")
	}
}

func TestDeleteRuchka(t *testing.T) {
	beforeServiceTest()
	ruchkas := []Ruchka{
		{
			Name:            "ruchka1",
			Uri:             "uri1",
			Method:          "method1",
			AllowedAccounts: []string{"acc1", "acc2"},
			AllowedGroups:   []string{"gr1", "gr2"},
		},
		{
			Name:            "ruchka2",
			Uri:             "uri2",
			Method:          "method2",
			AllowedAccounts: []string{"acc1", "acc2"},
			AllowedGroups:   []string{"gr1", "gr2"},
		},
		{
			Name:            "ruchka3",
			Uri:             "uri3",
			Method:          "method3",
			AllowedAccounts: []string{"acc1", "acc2"},
			AllowedGroups:   []string{"gr1", "gr2"},
		},
	}
	services := []ServiceInfo{
		{
			Name:            "service1",
			BaseUri:         "baseUri1",
			AllowedAccounts: []string{"acc1", "acc2"},
			AllowedGroups:   []string{"gr1", "gr2"},
			Ruchkas:         ruchkas,
		},
		{
			Name:            "service2",
			BaseUri:         "baseUri2",
			AllowedAccounts: []string{"acc1", "acc2"},
			AllowedGroups:   []string{"gr1", "gr2"},
			Ruchkas:         ruchkas,
		},
	}
	_, err := serviceCollection.InsertOne(context.Background(), services[0])
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = serviceCollection.InsertOne(context.Background(), services[1])
	err = DeleteRuchka("service1", ruchkas[1].Name)
	if err != nil {
		t.Errorf(err.Error())
	}
	var service ServiceInfo
	res := serviceCollection.FindOne(context.Background(), bson.D{{"Name", "service1"}})
	if res.Err() != nil {
		t.Errorf(res.Err().Error())
	}
	err = res.Decode(&service)
	if err != nil {
		t.Errorf(err.Error())
	}
	firstRuchka := 0
	thirdRuchka := 0
	for _, ruchka := range ruchkas {
		for _, dbRuchka := range service.Ruchkas {
			if ruchka.Name == "ruchka2" && ruchka.Name == dbRuchka.Name {
				t.Errorf("ruchka was not deleted")
			} else if reflect.DeepEqual(ruchka, dbRuchka) {
				if ruchka.Name == "ruchka1" {
					firstRuchka++
				} else if ruchka.Name == "ruchka3" {
					thirdRuchka++
				}
			}
		}
	}
	if firstRuchka != 1 && thirdRuchka != 1 {
		t.Errorf("wrong ruchka deleted")
	}
	if serviceCollection.FindOne(context.Background(), services[1]).Err() != nil {
		t.Errorf("wrong service modifyed")
	}
	err = DeleteRuchka(services[0].Name, ruchkas[1].Name)
	if err.Error() != "ruchka \""+ruchkas[1].Name+"\" not found" {
		t.Errorf("wrong error message for not existed ruchka deletion")
	}
}
