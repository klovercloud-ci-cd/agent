package service

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtService interface {
	ValidateToken(tokenString string) (bool, *jwt.Token)
}