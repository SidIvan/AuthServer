package route

import (
	"github.com/gorilla/mux"
)

func NewInternalRouter() *mux.Router {
	router := mux.NewRouter()
	return router
}
