package route

import (
	"AuthServer/internal/logic"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

func NewExternalRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/login", login)
	//router.HandleFunc("/auth", auth).Methods("GET")
	return router
}

func login(w http.ResponseWriter, r *http.Request) {
	authInfo := strings.Split(r.Header.Get("Oauth"), ":")
	if len(authInfo) != 3 {
		w.WriteHeader(http.StatusExpectationFailed)
		_, err := fmt.Fprintf(w, "Invalid args count, must be 3, was %d", len(authInfo))
		if err != nil {
			log.Print(err)
		}
		return
	}
	logic.Login(authInfo[0], authInfo[1], authInfo[2])
}

//func auth(w http.ResponseWriter, r *http.Request) {
//	oauthToken := r.Header.Get("Oauth")
//	fmt.Fprintf(w, logic.DecodeToken(oauthToken))
//}
