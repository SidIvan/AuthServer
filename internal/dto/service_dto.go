package dto

import (
	"AuthServer/internal/repo"
	"encoding/json"
)

type ServiceCreateIn struct {
	Name    string `json:"service name"`
	BaseUri string `json:"base uri"`
}

type AddRuchkaIn struct {
	ServiceName string      `json:"service name"`
	Ruchka      repo.Ruchka `json:"ruchka"`
}

type DeleteRuchkaIn struct {
	ServiceName string `json:"service name"`
	RuchkaName  string `json:"ruchka name"`
}

type ServiceDeleteIn struct {
	Name string `json:"service name"`
}

type ManagementAvailabilitySingleAccountIn struct {
	ServiceName string `json:"service name"`
	Login       string `json:"login"`
}

type ManagementAvailabilitySingleGroupIn struct {
	ServiceName string `json:"service name"`
	Group       string `json:"group"`
}

type ManagementAvailabilityToRuchkaSingleAccountIn struct {
	ServiceName string `json:"service name"`
	RuchkaName  string `json:"ruchka name"`
	Login       string `json:"login"`
}

type ManagementAvailabilityToRuchkaSingleGroupIn struct {
	ServiceName string `json:"service name"`
	RuchkaName  string `json:"ruchka name"`
	Group       string `json:"group"`
}

type SuccessCreateServiceOut struct {
	Result ResponseType `json:"type"`
	Id     string       `json:"id"`
}

//type ServiceInfo struct {
//	Name            string   `json:"name"`
//	BaseUri         string   `json:"base URI"`
//	AllowedAccounts []string `json:"allowed accounts"`
//	AllowedGroups   []string `json:"allowed groups"`
//	Ruchkas         []RuchkaInfo `json:"ruchkas"`
//}
//
//type RuchkaInfo struct {
//	Name            string `json:"name"`
//	Uri             string `json:"URI"`
//	Method          string `json:"method"`
//	AllowedAccounts []string `json:"allowed accounts"`
//	AllowedGroups   []string `json:"allowed groups"`
//}

type ServiceInfoSuccessOut struct {
	Result      ResponseType `json:"type"`
	ServiceInfo repo.Service `json:"service info"`
}

func NewSuccessCreateServiceOut(id string) SuccessCreateServiceOut {
	return SuccessCreateServiceOut{
		Result: OkR,
		Id:     id,
	}
}

func (r SuccessCreateServiceOut) RawBody() (ResponseType, []byte) {
	body, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return r.Result, body
}

func NewServiceInfoSuccessOut(service *repo.Service) ServiceInfoSuccessOut {
	return ServiceInfoSuccessOut{
		Result:      OkR,
		ServiceInfo: *service,
	}
}

func (r ServiceInfoSuccessOut) RawBody() (ResponseType, []byte) {
	body, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return r.Result, body
}
