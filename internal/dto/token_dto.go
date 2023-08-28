package dto

import "encoding/json"

type CreateRefreshIn struct {
	Login       string `json:"login"`
	ServiceName string `json:"service name"`
	RuchkaName  string `json:"ruchka name"`
}

type UpdateRefreshIn struct {
	TokenValue string `json:"token value"`
}

type TokenValueOut struct {
	Result     ResponseType `json:"result"`
	TokenValue string       `json:"token"`
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
