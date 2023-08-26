package route

import (
	"AuthServer/internal/dto"
	"AuthServer/internal/repo"
)

// TODO: test
func authorization(authIn dto.AuthIn, tokenValue string) dto.Response {
	token, err := repo.ParseToken(tokenValue)
	if err != nil {
		return dto.NewErrorOut("invalid token")
	}
	if !token.Valid {
		return dto.NewErrorOut("invalid token")
	}
	if !repo.IsAccess(tokenValue) {
		return dto.NewErrorOut("given token not access")
	}
	claims := token.Claims.(repo.Payload)
	if claims.Service != authIn.Service || claims.Login != authIn.Login || claims.Ruchka != authIn.Service {
		return dto.NewAuthorizationSuccessOut(false)
	}
	return dto.NewAuthorizationSuccessOut(true)
}
