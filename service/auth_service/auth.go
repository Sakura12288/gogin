package auth_service

import (
	"gogin/models"
	"gogin/pkg/util"
)

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Exist() (bool, error) {
	return models.CheckAuth(a.Username, a.Password)
}

func (a *Auth) GenerateRS256Token() (string, error) {
	return util.GenerateTokenUsingRS256(a.Username)
}
func (a *Auth) GenerateHS256Token() (string, error) {
	return util.GenerateTokenUsingHs256(a.Username, a.Password)
}
