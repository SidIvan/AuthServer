package dto

import "encoding/json"

type CreateRefreshIn struct {
	Login       string `json:"login"`
	ServiceName string `json:"service name"`
	RuchkaName  string `json:"ruchka name"`
}

type TokenValueIn struct {
	TokenValue string `json:"token value"`
}

type TokenValueOut struct {
	Result     ResponseType `json:"result"`
	TokenValue string       `json:"token"`
}

type IsValidOut struct {
	Result  ResponseType `json:"result"`
	IsValid bool         `json:"is valid"`
}

func (r TokenValueOut) RawBody() (ResponseType, []byte) {
	body, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return r.Result, body
}

func NewTokenValueOut(tokenValue string) TokenValueOut {
	return TokenValueOut{
		Result:     OkR,
		TokenValue: tokenValue,
	}
}

func (r IsValidOut) RawBody() (ResponseType, []byte) {
	body, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return r.Result, body
}

func NewIsValidOut(isValid bool) IsValidOut {
	return IsValidOut{
		Result:  OkR,
		IsValid: isValid,
	}
}
