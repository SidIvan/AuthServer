package route

import (
	"AuthServer/internal/dto"
	"AuthServer/internal/repo"
	"github.com/gorilla/mux"
	"net/http"
)

// TODO: test
func NewServiceRouter(router *mux.Router) {
	router = router.PathPrefix("/service").Subrouter()
	router.
		HandleFunc("/create", createServiceHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPost)
	router.
		HandleFunc("/info/{serviceName}", serviceInfoHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodGet)
	router.
		HandleFunc("/delete", deleteServiceHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodDelete)
	router.
		HandleFunc("/addRuchka", addRuchkaHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPut)
	router.
		HandleFunc("/deleteRuchka", deleteRuchkaHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPut)
	router.
		HandleFunc("/allowAcc", allowAccountHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPut)
	router.
		HandleFunc("/disallowAcc", disallowAccountHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPut)
	router.
		HandleFunc("/allowGroup", allowGroupHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPut)
	router.
		HandleFunc("/disallowGroup", disallowGroupHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPut)
	router.
		HandleFunc("/allowAccRuchka", allowAccountToRuchkaHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPut)
	router.
		HandleFunc("/disallowAccRuchka", disallowAccountToRuchkaHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPut)
	router.
		HandleFunc("/allowGroupRuchka", allowGroupToRuchkaHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPut)
	router.
		HandleFunc("/disallowGroupRuchka", disallowGroupToRuchkaHandler).
		Headers("Content-Type", "application/json", "Oauth", "").
		Methods(http.MethodPut)
}

// TODO: test
func createServiceHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.ServiceCreateIn
	if !checkAuthAndParseBody(w, r, &inInfo, "CreateService") {
		return
	}
	okOrErr(w, createService(inInfo))
}

// TODO: test
func createService(inInfo dto.ServiceCreateIn) dto.Response {
	id, err := repo.CreateService(inInfo.Name, inInfo.BaseUri)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewSuccessCreateServiceOut(id)
}

// TODO: test
func serviceInfoHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r, "ServiceInfo") {
		return
	}
	serviceName := mux.Vars(r)["serviceName"]
	okOrErr(w, getServiceInfo(serviceName))
}

// TODO: test
func getServiceInfo(sName string) dto.Response {
	service, err := repo.FindService(sName)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewServiceInfoSuccessOut(service)
}

// TODO: test
func deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.ServiceDeleteIn
	if !checkAuthAndParseBody(w, r, &inInfo, "DeleteService") {
		return
	}
	okOrErr(w, deleteService(inInfo))
}

// TODO: test
func deleteService(inInfo dto.ServiceDeleteIn) dto.Response {
	err := repo.DeleteService(inInfo.Name)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}

// TODO: test
func addRuchkaHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.AddRuchkaIn
	if !checkAuthAndParseBody(w, r, &inInfo, "AddRuchka") {
		return
	}
	okOrErr(w, addRuchka(&inInfo))
}

// TODO: test
func addRuchka(inInfo *dto.AddRuchkaIn) dto.Response {
	err := repo.AddRuchka(inInfo.ServiceName, inInfo.Ruchka)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}

// TODO: test
func deleteRuchkaHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.DeleteRuchkaIn
	if !checkAuthAndParseBody(w, r, &inInfo, "DeleteRuchka") {
		return
	}
	okOrErr(w, deleteRuchka(&inInfo))
}

// TODO: test
func deleteRuchka(inInfo *dto.DeleteRuchkaIn) dto.Response {
	err := repo.DeleteRuchka(inInfo.ServiceName, inInfo.RuchkaName)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}

// TODO: test
func allowAccountHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.ManagementAvailabilitySingleAccountIn
	if !checkAuthAndParseBody(w, r, &inInfo, "AllowAccount") {
		return
	}
	okOrErr(w, allowAccount(inInfo))
}

// TODO: test
func allowAccount(inInfo dto.ManagementAvailabilitySingleAccountIn) dto.Response {
	err := repo.AddAllowedAccount(inInfo.ServiceName, inInfo.Login)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}

// TODO: test
func disallowAccountHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.ManagementAvailabilitySingleAccountIn
	if !checkAuthAndParseBody(w, r, inInfo, "DisallowAcc") {
		return
	}
	okOrErr(w, disallowAccount(inInfo))
}

// TODO: test
func disallowAccount(inInfo dto.ManagementAvailabilitySingleAccountIn) dto.Response {
	err := repo.DeleteAllowedAccount(inInfo.ServiceName, inInfo.Login)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}

// TODO: test
func allowGroupHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.ManagementAvailabilitySingleGroupIn
	if !checkAuthAndParseBody(w, r, &inInfo, "AllowGroup") {
		return
	}
	okOrErr(w, allowGroup(inInfo))
}

// TODO: test
func allowGroup(inInfo dto.ManagementAvailabilitySingleGroupIn) dto.Response {
	err := repo.AddAllowedGroup(inInfo.ServiceName, inInfo.Group)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}

// TODO: test
func disallowGroupHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.ManagementAvailabilitySingleGroupIn
	if !checkAuthAndParseBody(w, r, &inInfo, "DisallowGroup") {
		return
	}
	okOrErr(w, disallowGroup(inInfo))
}

// TODO: test
func disallowGroup(inInfo dto.ManagementAvailabilitySingleGroupIn) dto.Response {
	err := repo.DeleteAllowedGroup(inInfo.ServiceName, inInfo.Group)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}

// TODO: test
func allowAccountToRuchkaHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.ManagementAvailabilityToRuchkaSingleAccountIn
	if !checkAuthAndParseBody(w, r, &inInfo, "AllowAccountToRuchka") {
		return
	}
	okOrErr(w, allowAccountToRuchka(inInfo))
}

// TODO: test
func allowAccountToRuchka(inInfo dto.ManagementAvailabilityToRuchkaSingleAccountIn) dto.Response {
	err := repo.AddRuchkaAllowedAccount(inInfo.ServiceName, inInfo.RuchkaName, inInfo.Login)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}

// TODO: test
func disallowAccountToRuchkaHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.ManagementAvailabilityToRuchkaSingleAccountIn
	if !checkAuthAndParseBody(w, r, &inInfo, "DisallowAccountToRuchka") {
		return
	}
	okOrErr(w, disallowAccountToRuchka(inInfo))
}

// TODO: test
func disallowAccountToRuchka(inInfo dto.ManagementAvailabilityToRuchkaSingleAccountIn) dto.Response {
	err := repo.DeleteRuchkaAllowedAccount(inInfo.ServiceName, inInfo.RuchkaName, inInfo.Login)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}

// TODO: test
func allowGroupToRuchkaHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.ManagementAvailabilityToRuchkaSingleGroupIn
	if !checkAuthAndParseBody(w, r, &inInfo, "AllowGroupToRuchka") {
		return
	}
	okOrErr(w, allowGroupToRuchka(inInfo))
}

// TODO: test
func allowGroupToRuchka(inInfo dto.ManagementAvailabilityToRuchkaSingleGroupIn) dto.Response {
	err := repo.AddRuchkaAllowedGroup(inInfo.ServiceName, inInfo.RuchkaName, inInfo.Group)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}

// TODO: test
func disallowGroupToRuchkaHandler(w http.ResponseWriter, r *http.Request) {
	var inInfo dto.ManagementAvailabilityToRuchkaSingleGroupIn
	if !checkAuthAndParseBody(w, r, &inInfo, "DisallowGroupToRuchka") {
		return
	}
	okOrErr(w, disallowGroupToRuchka(inInfo))
}

// TODO: test
func disallowGroupToRuchka(inInfo dto.ManagementAvailabilityToRuchkaSingleGroupIn) dto.Response {
	err := repo.DeleteRuchkaAllowedGroup(inInfo.ServiceName, inInfo.RuchkaName, inInfo.Group)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	return dto.NewOkOut()
}
