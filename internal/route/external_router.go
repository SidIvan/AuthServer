package route

import (
	"AuthServer/internal/dto"
	"AuthServer/internal/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func NewExternalRouter(router *mux.Router) {
	ThisServiceName = utils.PMan.Get("this_service_name").(string)
	router.
		HandleFunc("/registration", registrationHandler).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")
	router.
		HandleFunc("/authorization", authorizationHandler).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json",
			"Oauth", "")
}

// TODO: test
func authorizationHandler(w http.ResponseWriter, r *http.Request) {
	var authInfo dto.AuthIn
	if !checkAuthAndParseBody(w, r, &authInfo, "") {
		return
	}
	rType, body := authorization(authInfo).RawBody()
	if rType == dto.OkR {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadGateway)
	}
	numBytes, err := w.Write(body)
	if err != nil || numBytes != len(body) {
		log.Println("not full response sent")
	}
}

// TODO: test
func registrationHandler(w http.ResponseWriter, r *http.Request) {
	var regInfo dto.RegistrationIn
	body, err := io.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &regInfo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rType, body := registration(regInfo).RawBody()
	if rType == dto.OkR {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadGateway)
	}
	numBytes, err := w.Write(body)
	if err != nil || numBytes != len(body) {
		log.Println("not full response sent")
	}
}
