package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gogin/pkg/logging"
	"math/rand"
	"time"
)

const pri_key = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDxjb567s5NXCWlFk2JNFLkUN+E3Qm7OMRLxxDwJoSpbkYtxuEm
7UJEqx6hHHkY0lmCzkCXd8Ji3JKGaGflx+bFxyhj9mRWNFBuCSVXn9Cyzb2fC0kU
DSnKuWCDEVZQ5O5eQKzits1Oim3b+RbBYcAj8Dc4zQcwwPYuA+RgG/GAPQIDAQAB
AoGAMeXX7UkbcLuSQzICPk+CuAtEwQtwES6+zfCHPTSXvvA6qwYkSIhGYiz/HMTm
9wus1eqJSUDB9O4fjohOvha3QskOPQ+XilufcVWW6cE+z6CAwLePvUkyUSadx9wL
BUeuWoLTxH9aHpVORANd1LB63JN2rt6vQ8QZT1C7eb1gZrkCQQD5LTYQeA0R7AC7
NAhRg/D1bkvDlBtlQPei1i6sXrf/z7FvXt9vLtJgUut4Z2lUg8zlGqHTRM8nSOl5
9jc5eCynAkEA+CsYJ7euA8MlipmxyOZdB/43oM0L2bmoOI9tIbUfmD1hfyfT8klk
dmir3ga+YqP7tAcGvn5fkfD5xNaNKnEUewJAdlQEApogCsy6JCw3bw5rFQIFtKDW
yaSqdIelrnFki3SD3FF/ZXskqF14OLtTB7F3UazuADgC77LuPN6xpvbsrQJBAOOg
I3fKsoIg7L5EWx26rno2Yy/K46PA9ttqMt9IEsLBCjxne7AwQUWanIn6BYbUgnqO
N1Fi+KYUMgSqBrF3JyECQQC8pHILDRRZe96Y5PH2lvMhXdA9kDf2S3v2H9O0L4mp
surYRs9IyVCAsTHgbAHoVERPSE0M3d6eIS2tJNmMm3un
-----END RSA PRIVATE KEY-----
`
const pub_key = `-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBAPGNvnruzk1cJaUWTYk0UuRQ34TdCbs4xEvHEPAmhKluRi3G4SbtQkSr
HqEceRjSWYLOQJd3wmLckoZoZ+XH5sXHKGP2ZFY0UG4JJVef0LLNvZ8LSRQNKcq5
YIMRVlDk7l5ArOK2zU6Kbdv5FsFhwCPwNzjNBzDA9i4D5GAb8YA9AgMBAAE=
-----END RSA PUBLIC KEY-----
`

type RSClaim struct {
	Username  string
	RandomNum int
	jwt.RegisteredClaims
}

//生成rsa私钥

func parsePriKeyBytes(buf []byte) (*rsa.PrivateKey, error) {
	p := &pem.Block{}
	p, _ = pem.Decode(buf)
	if p == nil {
		return nil, errors.New("parse pri_key error")
	}
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}
func parsePubKeyBytes(buf []byte) (*rsa.PublicKey, error) {
	p := &pem.Block{}
	p, _ = pem.Decode(buf)
	if p == nil {
		return nil, errors.New("parse pub_key error")
	}
	return x509.ParsePKCS1PublicKey(p.Bytes)
}

func GenerateTokenUsingRS256(username string) (string, error) {
	num := rand.Intn(1 << 31)
	nowTime := time.Now()
	expireTime := time.Now().Add(2 * time.Hour)
	claim := RSClaim{
		username,
		num,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "TL",
		},
	}
	rsa_pri_key, err := parsePriKeyBytes([]byte(pri_key))
	if err != nil {
		fmt.Println(err.Error())
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claim).SignedString(rsa_pri_key)
	return token, err
}
func ParseTokenUsingRS256(token string) (*RSClaim, error) {
	tokenClaim, err := jwt.ParseWithClaims(token, &RSClaim{}, func(token *jwt.Token) (interface{}, error) {
		rsa_pub_key, err := parsePubKeyBytes([]byte(pub_key))
		if err != nil {
			logging.Info(err)
			return nil, err
		}
		return rsa_pub_key, nil
	})
	if claim, ok := tokenClaim.Claims.(*RSClaim); ok && tokenClaim.Valid {
		return claim, err
	}
	return nil, err
}
