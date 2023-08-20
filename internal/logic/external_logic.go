package logic

import (
	"AuthServer/internal/repo"
	"AuthServer/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"log"
)

func CreateAccount(login string, password string) {
	repo.InsertAccount(login, password)
}

func Login(login string, password string, service string) {
	repo.FindRefresh(login, service)
}

func verifyRefresh(signedToken string) (bool, *jwt.Token) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return utils.PMan.Get("HMAC_SECRET_KEY").([]byte), nil
	})
	if err != nil {
		log.Println("Can't parse token: ", err)
	}
	return token.Valid, token
}
