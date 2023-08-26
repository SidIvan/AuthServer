package dto

import "encoding/json"

type RegistrationIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegistrationSuccessOut struct {
	Result            ResponseType `json:"type"`
	RefreshTokenValue string       `json:"refresh token"`
	AccessTokenValue  string       `json:"access token"`
}

func NewRegistrationSuccessOut(refresh string, access string) RegistrationSuccessOut {
	return RegistrationSuccessOut{
		Result:            OkR,
		RefreshTokenValue: refresh,
		AccessTokenValue:  access,
	}
}

func (r RegistrationSuccessOut) RawBody() (ResponseType, []byte) {
	body, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return r.Result, body
}
