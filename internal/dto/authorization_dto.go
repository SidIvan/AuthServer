package dto

import "encoding/json"

type AuthIn struct {
	Login   string `json:"login"`
	Service string `json:"service"`
	Ruchka  string `json:"ruchka"`
}

type AuthOut struct {
	Result  ResponseType `json:"type"`
	Applied bool         `json:"applied"`
}

func NewAuthorizationSuccessOut(applied bool) AuthOut {
	return AuthOut{
		Result:  OkR,
		Applied: applied,
	}
}

func (r AuthOut) RawBody() (ResponseType, []byte) {
	body, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return r.Result, body
}
