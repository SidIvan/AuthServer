package route

import (
	"AuthServer/internal/dto"
	"AuthServer/internal/repo"
)

func isAllowed(tokenValue string, ruchkaName string) bool {
	token, err := repo.ParseToken(tokenValue)
	if err != nil || !token.Valid {
		return false
	}
	if !repo.IsAccess(tokenValue) {
		return false
	}
	claims := token.Claims.(*repo.Payload)
	if claims.Service != ThisServiceName {
		if claims.Ruchka != "" && claims.Ruchka != ruchkaName {
			return false
		}
	}
	if !repo.IsAccess(tokenValue) {
		return false
	}
	return true
}

// TODO: test
func authorization(authIn dto.AuthIn) dto.Response {
	token, err := repo.ParseToken(authIn.Token)
	if err != nil {
		return dto.NewErrorOut("invalid token")
	}
	if !token.Valid {
		return dto.NewErrorOut("invalid token")
	}
	if !repo.IsAccess(authIn.Token) {
		return dto.NewErrorOut("given token not access")
	}
	claims := token.Claims.(*repo.Payload)
	if claims.Service != authIn.Service || claims.Login != authIn.Login || claims.Ruchka != authIn.Ruchka {
		return dto.NewAuthorizationSuccessOut(false)
	}
	return dto.NewAuthorizationSuccessOut(true)
}
