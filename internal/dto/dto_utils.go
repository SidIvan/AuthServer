package dto

import "encoding/json"

type ResponseType string

const (
	OkR    ResponseType = "ok"
	ErrorR ResponseType = "error"
)

type Response interface {
	RawBody() (ResponseType, []byte)
}

type ErrorOut struct {
	Result  ResponseType `json:"type"`
	Message string       `json:"error message"`
}

func NewErrorOut(message string) ErrorOut {
	return ErrorOut{
		Result:  ErrorR,
		Message: message,
	}
}

func (r ErrorOut) RawBody() (ResponseType, []byte) {
	body, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return r.Result, body
}
