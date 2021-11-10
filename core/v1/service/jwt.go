package service

import (
	"github.com/dgrijalva/jwt-go"
)

// Jwt Jwt operations.
type Jwt interface {
	ValidateToken(tokenString string) (bool, *jwt.Token)
}
