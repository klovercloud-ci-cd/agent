package logic

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/klovercloud-ci/config"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"log"
)

type jwtService struct {
	Jwt v1.Jwt
}

func (j jwtService) ValidateToken(tokenString string) (bool, *jwt.Token) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return (j.Jwt.PublicKey), nil
	})
	if err!=nil{
		log.Print("[ERROR]: Token is invalid! ",err.Error())
		return false,nil
	}
	return true,token


}


func getPublicKey() *rsa.PublicKey {
	block, _ := pem.Decode([]byte(config.Publickey))
	publicKeyImported, err := x509.ParsePKCS1PublicKey(block.Bytes)

	if err != nil {
		log.Print(err.Error())
		panic(err)
	}
	return publicKeyImported
}

func NewJwtService() service.JwtService {
	return jwtService{
		Jwt: v1.Jwt{
			PublicKey:  getPublicKey(),
		},
	}
}
