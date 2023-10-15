package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret []byte

type HSClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateTokenUsingHs256(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claim := HSClaims{
		username,
		password,
		jwt.StandardClaims{
			Issuer:    "TL",
			Subject:   username,
			IssuedAt:  nowTime.Unix(),
			ExpiresAt: expireTime.Unix(),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	fmt.Println(tokenClaims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}
func ParseToken(token string) (*HSClaims, error) {
	tokenClaim, err := jwt.ParseWithClaims(token, &HSClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaim != nil {
		if claim, ok := tokenClaim.Claims.(*HSClaims); ok && claim != nil {
			return claim, nil
		}
	}
	return nil, err
}
