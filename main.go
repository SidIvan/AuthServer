package main

import (
	"AuthServer/internal/repo"
	"AuthServer/internal/route"
	"AuthServer/internal/utils"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	utils.PMan = utils.NewPman("application.properties")
	ctx := context.Background()
	repo.ConnectToMongo(ctx, "mongodb://"+utils.PMan.Get("mongo_host").(string)+":"+utils.PMan.Get("mongo_port").(string), utils.PMan.Get("mongo_db_name").(string))
	repo.DropDb()
	router := mux.NewRouter()
	route.NewExternalRouter(router)
	route.NewServiceRouter(router)
	route.TokenRouterConfig(router)
	http.Handle("/", router)
	//serviceRouter := route.NewServiceRouter()
	//http.Handle("/service", serviceRouter)
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
		Name:            "ServiceInfo",
		Uri:             "/service/info/{serviceName}",
		Method:          http.MethodGet,
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
		Name:            "AddRuchka",
		Uri:             "/addRuchka",
		Method:          http.MethodPut,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})

	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "DeleteRuchka",
		Uri:             "/deleteRuchka",
		Method:          http.MethodPut,
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
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "CreateRefresh",
		Uri:             "/refresh/create",
		Method:          http.MethodPost,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "UpdateRefresh",
		Uri:             "/refresh/update",
		Method:          http.MethodPatch,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
	repo.AddRuchka(route.ThisServiceName, repo.Ruchka{
		Name:            "ValidateRefresh",
		Uri:             "/refresh/validate/{token}",
		Method:          http.MethodGet,
		AllowedAccounts: []string{"DrLivesey"},
		AllowedGroups:   nil,
	})
}
