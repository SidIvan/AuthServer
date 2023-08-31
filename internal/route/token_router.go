package route

import (
	"AuthServer/internal/dto"
	"AuthServer/internal/repo"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
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

// TODO: test
func updateRefreshHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.TokenValueIn
	if !checkAuthAndParseBody(w, r, &inInfo, "UpdateRefresh") {
		return
	}
	okOrErr(w, updateRefresh(&inInfo))
}

// TODO: test
func updateRefresh(inInfo *dto.TokenValueIn) dto.Response {
	token, err := repo.ParseToken(inInfo.TokenValue)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	tokenValue, err := repo.CreateRefresh(repo.Payload{
		Login:   claims["Login"].(string),
		Service: claims["Service"].(string),
		Ruchka:  claims["Ruchka"].(string),
	})
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewTokenValueOut(tokenValue)
}

// TODO: test
func validateRefreshHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.TokenValueIn
	if !checkAuthAndParseBody(w, r, &inInfo, "ValidateRefresh") {
		return
	}
	okOrErr(w, validateRefresh(&inInfo))
}

// TODO: test
func validateRefresh(inInfo *dto.TokenValueIn) dto.Response {
	return dto.NewIsValidOut(repo.IsRefresh(inInfo.TokenValue))

}
