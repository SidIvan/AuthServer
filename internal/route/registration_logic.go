package route

import (
	"AuthServer/internal/dto"
	"AuthServer/internal/repo"
)

var thisServiceName string

// TODO: test
func registration(regInfo dto.RegistrationIn) dto.Response {
	payload := repo.Payload{
		Login:   regInfo.Login,
		Service: thisServiceName,
	}
	_, err := repo.CreateAccount(regInfo.Login, regInfo.Password)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	refreshTokenValue, err := repo.CreateRefresh(payload)
	if err != nil {
		return dto.NewErrorOut(err.Error())
	}
	accessTokenValue, err := repo.CreateAccess(payload)
	if err != nil {
		errDelRef := repo.DeleteRefresh(refreshTokenValue)
		if errDelRef == nil {
			return dto.NewErrorOut(err.Error())
		}
		panic("invalid token in DB")
	}
	return dto.NewRegistrationSuccessOut(refreshTokenValue, accessTokenValue)
}
