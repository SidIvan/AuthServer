package main

import (
	"AuthServer/internal/repo"
	"AuthServer/internal/route"
	"AuthServer/internal/utils"
	"context"
	"fmt"
	"net/http"
)

func main() {
	utils.PMan = utils.NewPman("application.properties")
	ctx := context.Background()
	repo.ConnectToMongo(ctx, "mongodb://"+utils.PMan.Get("mongo_host").(string)+":"+utils.PMan.Get("mongo_port").(string), utils.PMan.Get("mongo_db_name").(string))
	repo.DropDb()
	extRouter := route.NewExternalRouter()
	http.Handle("/", extRouter)
	serviceRouter := route.NewServiceRouter()
	http.Handle("/service", serviceRouter)
	repo.CreateAccount("DrLivesey", "Rum")
	repo.CreateService(utils.PMan.Get("this_service_name").(string), "")
	fmt.Println(repo.CreateAccess(repo.Payload{
		Login:     "DrLivesey",
		Service:   route.ThisServiceName,
		Ruchka:    "",
		MapClaims: nil,
	}))
	defaultRuchkas()
	http.ListenAndServe(":8181", nil)
}

func defaultRuchkas() {
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "CreateService",
		Uri:             "/service/create",
		Method:          http.MethodPost,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "DeleteService",
		Uri:             "/service/delete",
		Method:          http.MethodDelete,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "AllowAccount",
		Uri:             "/service/allowAcc",
		Method:          http.MethodPut,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "DisallowAcc",
		Uri:             "/service/allowAcc",
		Method:          http.MethodPut,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "AllowGroup",
		Uri:             "/service/allowGroup",
		Method:          http.MethodPut,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "DisallowGroup",
		Uri:             "/service/disallowGroup",
		Method:          http.MethodPut,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "AllowAccountToRuchka",
		Uri:             "/service/allowAccRuchka",
		Method:          http.MethodPut,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "DisallowAccountToRuchka",
		Uri:             "/service/disallowAccRuchka",
		Method:          http.MethodPut,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "AllowGroupToRuchka",
		Uri:             "/service/allowGroupRuchka",
		Method:          http.MethodPut,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "DisallowGroupToRuchka",
		Uri:             "/service/disallowGroupRuchka",
		Method:          http.MethodPut,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
}
