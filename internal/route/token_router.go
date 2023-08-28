package route

import (
	"AuthServer/internal/dto"
	"AuthServer/internal/repo"
	"github.com/gorilla/mux"
	"net/http"
)

func TokenRouterConfig(router *mux.Router) {
	router = router.PathPrefix("/token").Subrouter()
	router.
		HandleFunc("/create", createRefreshHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPost)
	router.
		HandleFunc("/update", updateRefreshHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPatch)
	router.
		HandleFunc("/validate", validateRefreshHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPost)
}

// TODO: test
func createRefreshHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.CreateRefreshIn
	if !checkAuthAndParseBody(w, r, &inInfo, "CreateRefresh") {
		return
	}
	okOrErr(w, createRefresh(&inInfo))
}

// TODO: test
func createRefresh(inInfo *dto.CreateRefreshIn) dto.Response {
	tokenValue, err := repo.CreateRefresh(repo.Payload{
		Login:   inInfo.Login,
		Service: inInfo.ServiceName,
		Ruchka:  inInfo.RuchkaName,
	})
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewTokenValueOut(tokenValue)
}
