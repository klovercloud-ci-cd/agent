package v1

import (
	"crypto/rsa"
)

// Jwt Struct of Jwt keys
type Jwt struct {
	PublicKey *rsa.PublicKey
}
