package dto

import "encoding/json"

type ServiceCreateIn struct {
	Name    string `json:"service name"`
	BaseUri string `json:"base uri"`
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
