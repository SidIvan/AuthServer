package route

import (
	"AuthServer/internal/dto"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// TODO: test
func checkAuth(w http.ResponseWriter, r *http.Request, ruchkaName string) bool {
	tokenValue := r.Header.Get("Oauth")
	if !isAllowed(tokenValue, ruchkaName) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func parseBody(w http.ResponseWriter, r *http.Request, bodyHandler interface{}) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	err = json.Unmarshal(body, bodyHandler)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}

// TODO: test
func checkAuthAndParseBody(w http.ResponseWriter, r *http.Request, bodyHandler interface{}, ruchkaName string) bool {
	if !checkAuth(w, r, ruchkaName) {
		return false
	}
	return parseBody(w, r, bodyHandler)
}

// TODO: test
func okOrErr(w http.ResponseWriter, response dto.Response) {
	responseType, body := response.RawBody()
	if responseType == dto.OkR {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadGateway)
	}
	numBytes, err := w.Write(body)
	if err != nil || numBytes != len(body) {
		log.Println("not full response sent")
	}
}
