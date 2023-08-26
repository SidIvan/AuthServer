package main

import (
	"AuthServer/internal/repo"
	"AuthServer/internal/route"
	"AuthServer/internal/utils"
	"context"
	"net/http"
)

func main() {
	utils.PMan = utils.NewPman("application.properties")
	ctx := context.Background()
	repo.ConnectToMongo(ctx, "mongodb://"+utils.PMan.Get("mongo_host").(string)+":"+utils.PMan.Get("mongo_port").(string), utils.PMan.Get("mongo_db_name").(string))
	repo.DropDb()
	extRouter := route.NewExternalRouter()
	http.Handle("/", extRouter)
	repo.CreateService(utils.PMan.Get("this_service_name").(string), "")
	http.ListenAndServe(":8181", nil)
}
